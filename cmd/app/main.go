package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/sync/errgroup"

	v1 "homework/internal/api/v1"
	"homework/internal/config"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/specs"
)

func main() {
	var (
		err         error
		ctx, cancel = signal.NotifyContext(
			context.Background(),
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
	)
	defer cancel()

	cfg, err := config.InitConfig(os.Args)
	if err != nil {
		log.Fatal("get config: ", err.Error())
		return
	}

	// инициализация пакета/драйвера БД
	db, err := pgxpool.Connect(ctx, cfg.DB.Postgresql)
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// инициализация хранилищ
	storageRegistry := storage.NewStorageRegistry(cfg, db)

	// инициализация сервисов
	serviceRegistry := service.NewServiceRegistry(cfg, storageRegistry)

	// инициализация хэндлеров
	apiServer := v1.NewAPIServer(serviceRegistry)

	// запуск HTTP сервера
	log.Println("start HTTP server")
	err = startHTTPServer(ctx, cfg, apiServer)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
	log.Println("exit")

}

func startHTTPServer(
	ctx context.Context,
	cfg *config.Config,
	apiServer specs.ServerInterface,
	middlewares ...specs.MiddlewareFunc,
) error {

	handler := specs.HandlerWithOptions(apiServer, specs.ChiServerOptions{
		BaseURL:     cfg.BasePath,
		Middlewares: middlewares,
	})

	router := chi.NewRouter()
	router.Use(commonMiddleware)
	router.Handle("/*", handler)

	httpServer := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		return httpServer.Shutdown(ctx)
	})

	return group.Wait()
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
