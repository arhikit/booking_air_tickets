package flights

import (
	"github.com/google/uuid"
	"time"
)

type Airline struct {
	Id   uuid.UUID
	Name string
}

type Aircraft struct {
	Id      uuid.UUID
	Airline Airline
	Name    string
}

type City struct {
	Id   uuid.UUID
	Name string
}

type Airport struct {
	Id   uuid.UUID
	City City
	Name string
}

type ClassSeats struct {
	Id         uuid.UUID
	Aircraft   Aircraft
	Name       string
	CountSeats int
	Width      int
	Pitch      int
	CountInRow int
}

type Seat struct {
	Id         uuid.UUID
	ClassSeats ClassSeats
	Number     string
}

type FlightPrice struct {
	ClassSeats       ClassSeats
	CountVacantSeats int
	PriceTicket      int
}

type Flight struct {
	Id                     uuid.UUID
	Name                   string
	Aircraft               Aircraft
	DepartureAirport       Airport
	ArrivalAirport         Airport
	DepartureDate          time.Time
	Duration               time.Duration
	PricesTickets          []FlightPrice
	PriceAdditionalBaggage int
	PriceSeatSelection     int
	IsInternational        bool
	BaggageIncluded        bool
	PetAllowed             bool
}

// структура, содержащая параметры метода GetFlights
type ParamsGetFlights struct {
	DepartureCityId uuid.UUID
	ArrivalCityId   uuid.UUID
	DepartureDate   time.Time
}

// структура, используемая как вывода результата метода GetFlightVacantSeats,
// а также для проверок при создании и регистрации билета
type VacantSeats struct {
	ClassSeatsId     uuid.UUID
	ClassSeatsName   string
	CountVacantSeats int
	Seats            []Seat
}
