package v1

import (
	"homework/internal/service"
	"homework/specs"
)

// Быстрая проверка актуальности текущего интерфейса сервера.
var _ specs.ServerInterface = &apiServer{}

type apiServer struct {
	serviceRegistry *service.Services
}

func NewAPIServer(serviceRegistry *service.Services) specs.ServerInterface {
	return &apiServer{
		serviceRegistry: serviceRegistry,
	}
}
