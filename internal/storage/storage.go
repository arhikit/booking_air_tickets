package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"homework/internal/config"
	flightsStorage "homework/internal/storage/flights"
	ticketsStorage "homework/internal/storage/tickets"
	usersStorage "homework/internal/storage/users"
)

type Storages struct {
	Flight flightsStorage.FlightsStorage
	Ticket ticketsStorage.TicketsStorage
	User   usersStorage.UsersStorage
}

func NewStorageRegistry(cfg *config.Config, db *pgxpool.Pool) *Storages {

	flight := flightsStorage.NewFlightsStorage(db)
	ticket := ticketsStorage.NewTicketsStorage(db)
	user := usersStorage.NewUsersStorage(db)

	return &Storages{
		Flight: flight,
		Ticket: ticket,
		User:   user,
	}
}
