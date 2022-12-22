package tickets

import (
	"time"

	"github.com/google/uuid"

	flightsDomain "homework/internal/domain/flights"
	usersDomain "homework/internal/domain/users"
)

type Status struct {
	Id        int
	Name      string
	Timestamp time.Time
}

type Passenger struct {
	Id                    uuid.UUID
	User                  usersDomain.User
	NamePassenger         string
	IdentityDataPassenger string
}

type Ticket struct {
	Id                     uuid.UUID
	Status                 Status
	Flight                 flightsDomain.Flight
	User                   usersDomain.User
	Passenger              Passenger
	ClassSeats             flightsDomain.ClassSeats
	Seat                   *flightsDomain.Seat
	CountAdditionalBaggage int
	Price                  int
	PaidWithBonuses        int
	AccruedBonuses         int
}

// структуры, содержащие параметры методов:

type ParamsCreateTicket struct {
	StatusTimestamp        time.Time
	FlightId               uuid.UUID
	UserId                 uuid.UUID
	PassengerId            *uuid.UUID
	ParamsCreatePassenger  *ParamsCreatePassenger
	ClassSeatsId           uuid.UUID
	SeatId                 *uuid.UUID
	CountAdditionalBaggage int
	Price                  int
}

type ParamsCreatePassenger struct {
	NamePassenger         string
	IdentityDataPassenger string
}

type ParamsPayForTicket struct {
	StatusTimestamp time.Time
	TicketId        uuid.UUID
	UserId          uuid.UUID
	UserBalanceInit bool
	Price           int
	PaidWithBonuses int
	AccruedBonuses  int
}

type ParamsRefundTicket struct {
	StatusTimestamp time.Time
	TicketId        uuid.UUID
	UserId          uuid.UUID
	Price           int
}

type ParamsRegisterTicket struct {
	StatusTimestamp time.Time
	TicketId        uuid.UUID
	UserId          uuid.UUID
	SeatId          *uuid.UUID
	AccruedBonuses  int
}
