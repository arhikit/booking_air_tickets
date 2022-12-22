package users

import "github.com/google/uuid"

type UserBalance struct {
	SumPurchases int
	SumBonuses   int
}

type User struct {
	Id      uuid.UUID
	Name    string
	Email   string
	Balance *UserBalance
}
