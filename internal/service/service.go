package service

import (
	"homework/internal/config"
	flightsService "homework/internal/service/flights"
	ticketsService "homework/internal/service/tickets"
	usersService "homework/internal/service/users"
	storage "homework/internal/storage"
)

type Services struct {
	Flight flightsService.FlightsService
	Ticket ticketsService.TicketsService
	User   usersService.UsersService
}

func NewServiceRegistry(cfg *config.Config, Storages *storage.Storages) *Services {

	flight := flightsService.NewFlightsService(
		Storages.Flight)
	ticket := ticketsService.NewTicketsService(
		Storages.Ticket,
		Storages.Flight,
		Storages.User,
	)
	user := usersService.NewUsersService(
		Storages.User)

	return &Services{
		Flight: flight,
		Ticket: ticket,
		User:   user,
	}
}
