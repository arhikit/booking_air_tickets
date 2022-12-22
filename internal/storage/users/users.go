package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	usersDomain "homework/internal/domain/users"
	terr "homework/internal/util/terr"
)

type UsersStorage interface {
	GetAccruedBonuses(ctx context.Context, userId uuid.UUID, ticketPrice int) (int, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error)
}

type storage struct {
	db *pgxpool.Pool
}

func (s storage) GetAccruedBonuses(ctx context.Context, userId uuid.UUID, ticketPrice int) (int, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return 0, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`SELECT 
				bonus_calc_scale.percent
     		FROM users_balance
      			INNER JOIN bonus_calc_scale
     				ON users_balance.sum_purchases BETWEEN bonus_calc_scale.sum_purchases_to AND bonus_calc_scale.sum_purchases_from
			WHERE users_balance.user_id = $1;`,
		userId.String())

	var percent int
	err = row.Scan(
		&percent,
	)

	if err != nil && err != pgx.ErrNoRows {
		return 0, terr.SQLDatabaseError(err)
	}

	accruedBonuses := int((percent * ticketPrice) / 100)
	return accruedBonuses, nil
}

func (s storage) GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`
	 		SELECT 
				users.id,
				users.name,
				users.email,
				CASE
					WHEN users_balance.user_id IS NOT NULL
						THEN true
					ELSE false
				END users_balance_exists,
				CASE
					WHEN users_balance.user_id IS NOT NULL
						THEN users_balance.sum_purchases
					ELSE 0
				END sum_purchases,
				CASE
					WHEN users_balance.user_id IS NOT NULL
						THEN users_balance.sum_bonuses
					ELSE 0
				END sum_bonuses
	 		FROM users
				LEFT JOIN users_balance
					ON users.id = users_balance.user_id
			WHERE users.id = $1`,
		userId.String())

	var user usersDomain.User
	var userBalanceExists bool
	var balance usersDomain.UserBalance
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&userBalanceExists,
		&balance.SumPurchases,
		&balance.SumBonuses,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found user (id %s)", userId))

		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}

	if userBalanceExists {
		user.Balance = &balance
	}
	return &user, nil
}

func NewUsersStorage(db *pgxpool.Pool) UsersStorage {
	return &storage{db: db}
}
