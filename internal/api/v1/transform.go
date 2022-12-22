package v1

import (
	"errors"
	uuid "github.com/google/uuid"
	"time"

	flightsDomain "homework/internal/domain/flights"
	ticketsDomain "homework/internal/domain/tickets"
	usersDomain "homework/internal/domain/users"
	terr "homework/internal/util/terr"
	specs "homework/specs"
)

func convertStringToUuid(idString string) (uuid.UUID, error) {

	var id uuid.UUID

	if idString == "" {
		return id, errors.New("empty uuid")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return id, err
	}

	return id, nil
}

func transformParamsGetFlights(paramsFlightsSpecs *specs.GetFlightsParams) (*flightsDomain.ParamsGetFlights, error) {

	departureCityId, err := convertStringToUuid(paramsFlightsSpecs.DepartureCityId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_DEPARTURE_CITY_UUID", err.Error())
	}

	arrivalCityId, err := convertStringToUuid(paramsFlightsSpecs.ArrivalCityId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_ARRIVAL_CITY_UUID", err.Error())
	}

	var paramsGetFlights flightsDomain.ParamsGetFlights
	paramsGetFlights.DepartureCityId = departureCityId
	paramsGetFlights.ArrivalCityId = arrivalCityId
	paramsGetFlights.DepartureDate = paramsFlightsSpecs.DepartureDate.Time

	return &paramsGetFlights, nil
}

func transformParamsCreateTicket(paramsCreateTicketSpecs *specs.ParamsCreateTicket) (*ticketsDomain.ParamsCreateTicket, error) {

	flightId, err := convertStringToUuid(paramsCreateTicketSpecs.FlightId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_FLIGHT_UUID", err.Error())
	}

	userId, err := convertStringToUuid(paramsCreateTicketSpecs.UserId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_USER_UUID", err.Error())
	}

	// если передается PassengerId, значит используется уже существующий пассажир и нового создавать не надо
	passengerExists := paramsCreateTicketSpecs.PassengerId != nil
	var passengerId uuid.UUID
	if passengerExists {
		passengerId, err = convertStringToUuid(*paramsCreateTicketSpecs.PassengerId)
		if err != nil {
			return nil, terr.BadRequest("INVALID_PASSENGER_UUID", err.Error())
		}
	} else {
		if paramsCreateTicketSpecs.NamePassenger == nil || *paramsCreateTicketSpecs.NamePassenger == "" {
			return nil, terr.BadRequest("INVALID_NAME_PASSENGER", "empty name passenger")
		}
		if paramsCreateTicketSpecs.IdentityDataPassenger == nil || *paramsCreateTicketSpecs.IdentityDataPassenger == "" {
			return nil, terr.BadRequest("INVALID_IDENTITY_DATA_PASSENGER", "empty identity data passenger")
		}
	}

	classSeatsId, err := convertStringToUuid(paramsCreateTicketSpecs.ClassSeatsId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_CLASS_SEAT_UUID", err.Error())
	}

	// если передается SeatId, значит пассажир уже выбрал определенное место
	isSeatAssigned := paramsCreateTicketSpecs.SeatId != nil
	var seatId uuid.UUID
	if isSeatAssigned {
		seatId, err = convertStringToUuid(*paramsCreateTicketSpecs.SeatId)
		if err != nil {
			return nil, terr.BadRequest("INVALID_SEAT_UUID", err.Error())
		}
	}

	var paramsCreateTicket ticketsDomain.ParamsCreateTicket
	paramsCreateTicket.StatusTimestamp = time.Now()
	paramsCreateTicket.FlightId = flightId
	paramsCreateTicket.UserId = userId
	paramsCreateTicket.ClassSeatsId = classSeatsId
	paramsCreateTicket.CountAdditionalBaggage = paramsCreateTicketSpecs.CountAdditionalBaggage

	if passengerExists {
		paramsCreateTicket.PassengerId = &passengerId
	} else {
		paramsCreateTicket.ParamsCreatePassenger = &ticketsDomain.ParamsCreatePassenger{
			NamePassenger:         *paramsCreateTicketSpecs.NamePassenger,
			IdentityDataPassenger: *paramsCreateTicketSpecs.IdentityDataPassenger,
		}
	}

	if isSeatAssigned {
		paramsCreateTicket.SeatId = &seatId
	}

	return &paramsCreateTicket, nil
}

func transformParamsPayForTicket(paramsPayForTicketSpecs *specs.ParamsPayForTicket) (*ticketsDomain.ParamsPayForTicket, error) {

	ticketId, err := convertStringToUuid(paramsPayForTicketSpecs.TicketId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_TICKET_UUID", err.Error())
	}

	userId, err := convertStringToUuid(paramsPayForTicketSpecs.UserId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_USER_UUID", err.Error())
	}

	if paramsPayForTicketSpecs.PaidWithBonuses < 0 {
		return nil, terr.BadRequest("INVALID_SUM_BONUSES", "Bonuses sum is a positive number")
	}

	var paramsPayForTicket ticketsDomain.ParamsPayForTicket
	paramsPayForTicket.StatusTimestamp = time.Now()
	paramsPayForTicket.TicketId = ticketId
	paramsPayForTicket.UserId = userId
	paramsPayForTicket.PaidWithBonuses = paramsPayForTicketSpecs.PaidWithBonuses

	return &paramsPayForTicket, nil
}

func transformParamsRefundTicket(paramsRefundTicketSpecs *specs.ParamsRefundTicket) (*ticketsDomain.ParamsRefundTicket, error) {

	ticketId, err := convertStringToUuid(paramsRefundTicketSpecs.TicketId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_TICKET_UUID", err.Error())
	}

	userId, err := convertStringToUuid(paramsRefundTicketSpecs.UserId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_USER_UUID", err.Error())
	}

	var paramsRefundTicket ticketsDomain.ParamsRefundTicket
	paramsRefundTicket.StatusTimestamp = time.Now()
	paramsRefundTicket.TicketId = ticketId
	paramsRefundTicket.UserId = userId

	return &paramsRefundTicket, nil
}

func transformParamsRegisterTicket(paramsRegisterTicketSpecs *specs.ParamsRegisterTicket) (*ticketsDomain.ParamsRegisterTicket, error) {

	ticketId, err := convertStringToUuid(paramsRegisterTicketSpecs.TicketId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_TICKET_UUID", err.Error())
	}

	userId, err := convertStringToUuid(paramsRegisterTicketSpecs.UserId)
	if err != nil {
		return nil, terr.BadRequest("INVALID_USER_UUID", err.Error())
	}

	// SeatId передается для тех билетов, у которых не было выбрано место на этапе покупки билета
	isSeatAssigned := paramsRegisterTicketSpecs.SeatId != nil
	var seatId uuid.UUID
	if isSeatAssigned {
		seatId, err = convertStringToUuid(*paramsRegisterTicketSpecs.SeatId)
		if err != nil {
			return nil, terr.BadRequest("INVALID_SEAT_UUID", err.Error())
		}
	}

	var paramsRegisterTicket ticketsDomain.ParamsRegisterTicket
	paramsRegisterTicket.StatusTimestamp = time.Now()
	paramsRegisterTicket.TicketId = ticketId
	paramsRegisterTicket.UserId = userId

	if isSeatAssigned {
		paramsRegisterTicket.SeatId = &seatId
	}
	return &paramsRegisterTicket, nil
}

func transformFlight(flight *flightsDomain.Flight) *specs.Flight {

	var flightSpec specs.Flight

	flightSpec.Id = flight.Id.String()
	flightSpec.Name = flight.Name

	flightSpec.Airline.AirlineId = flight.Aircraft.Airline.Id.String()
	flightSpec.Airline.AirlineName = flight.Aircraft.Airline.Name
	flightSpec.Airline.AircraftId = flight.Aircraft.Id.String()
	flightSpec.Airline.AircraftName = flight.Aircraft.Name

	flightSpec.Departure.CityId = flight.DepartureAirport.City.Id.String()
	flightSpec.Departure.CityName = flight.DepartureAirport.City.Name
	flightSpec.Departure.AirportId = flight.DepartureAirport.Id.String()
	flightSpec.Departure.AirportName = flight.DepartureAirport.Name

	flightSpec.Arrival.CityId = flight.ArrivalAirport.City.Id.String()
	flightSpec.Arrival.CityName = flight.ArrivalAirport.City.Name
	flightSpec.Arrival.AirportId = flight.ArrivalAirport.Id.String()
	flightSpec.Arrival.AirportName = flight.ArrivalAirport.Name

	flightSpec.Date.Departure = flight.DepartureDate
	flightSpec.Date.Arrival = flight.DepartureDate.Add(flight.Duration)
	flightSpec.Date.Duration = int(flight.Duration / time.Minute)

	PricesTickets := make([]specs.FlightPrice, len(flight.PricesTickets))
	for i, flightPrice := range flight.PricesTickets {
		PricesTickets[i].ClassSeatsId = flightPrice.ClassSeats.Id.String()
		PricesTickets[i].ClassSeatsName = flightPrice.ClassSeats.Name
		PricesTickets[i].CountVacantSeats = flightPrice.CountVacantSeats
		PricesTickets[i].PriceTicket = flightPrice.PriceTicket
	}
	flightSpec.PricesTickets = PricesTickets

	flightSpec.PriceAdditionalBaggage = flight.PriceAdditionalBaggage
	flightSpec.PriceSeatSelection = flight.PriceSeatSelection

	flightSpec.IsInternational = flight.IsInternational
	flightSpec.BaggageIncluded = flight.BaggageIncluded
	flightSpec.PetAllowed = flight.PetAllowed

	return &flightSpec
}

func transformVacantSeats(vacantSeats *flightsDomain.VacantSeats) *specs.VacantSeats {

	var vacantSeatsSpec specs.VacantSeats

	vacantSeatsSpec.ClassSeatsId = vacantSeats.ClassSeatsId.String()
	vacantSeatsSpec.ClassSeatsName = vacantSeats.ClassSeatsName
	vacantSeatsSpec.CountVacantSeats = vacantSeats.CountVacantSeats

	Seats := make([]specs.Seat, len(vacantSeats.Seats))
	for i, seat := range vacantSeats.Seats {
		Seats[i].Id = seat.Id.String()
		Seats[i].Number = seat.Number
	}
	vacantSeatsSpec.Seats = Seats
	return &vacantSeatsSpec
}

func transformTicket(ticket *ticketsDomain.Ticket) *specs.Ticket {

	var ticketSpecs specs.Ticket

	ticketSpecs.Id = ticket.Id.String()

	ticketSpecs.Status.Name = ticket.Status.Name
	ticketSpecs.Status.Timestamp = ticket.Status.Timestamp

	ticketSpecs.Flight.Id = ticket.Flight.Id.String()
	ticketSpecs.Flight.Name = ticket.Flight.Name
	ticketSpecs.Flight.Airline = ticket.Flight.Aircraft.Airline.Name
	ticketSpecs.Flight.Aircraft = ticket.Flight.Aircraft.Name
	ticketSpecs.Flight.DepartureCity = ticket.Flight.DepartureAirport.City.Name
	ticketSpecs.Flight.DepartureAirport = ticket.Flight.DepartureAirport.Name
	ticketSpecs.Flight.DepartureDate = ticket.Flight.DepartureDate
	ticketSpecs.Flight.ArrivalCity = ticket.Flight.ArrivalAirport.City.Name
	ticketSpecs.Flight.ArrivalAirport = ticket.Flight.ArrivalAirport.Name
	ticketSpecs.Flight.ArrivalDate = ticket.Flight.DepartureDate.Add(ticket.Flight.Duration)
	ticketSpecs.Flight.Duration = int(ticket.Flight.Duration / time.Minute)

	ticketSpecs.User.Id = ticket.User.Id.String()
	ticketSpecs.User.Name = ticket.User.Name

	ticketSpecs.Passenger.Id = ticket.Passenger.Id.String()
	ticketSpecs.Passenger.Name = ticket.Passenger.NamePassenger
	ticketSpecs.Passenger.IdentityData = ticket.Passenger.IdentityDataPassenger

	ticketSpecs.Seat.ClassSeatsId = ticket.ClassSeats.Id.String()
	ticketSpecs.Seat.ClassSeatsName = ticket.ClassSeats.Name
	if ticket.Seat != nil {
		seatId := ticket.Seat.Id.String()
		ticketSpecs.Seat.SeatId = &seatId
		ticketSpecs.Seat.SeatNumber = &ticket.Seat.Number
	}

	ticketSpecs.СountAdditionalBaggage = ticket.CountAdditionalBaggage
	ticketSpecs.Price = ticket.Price
	ticketSpecs.PaidWithBonuses = ticket.PaidWithBonuses
	ticketSpecs.AccruedBonuses = ticket.AccruedBonuses

	return &ticketSpecs
}

func transformUser(user *usersDomain.User) *specs.User {

	var userSpecs specs.User
	userSpecs.Id = user.Id.String()
	userSpecs.Name = user.Name
	userSpecs.Email = user.Email

	if user.Balance != nil {
		userSpecs.Balance.SumPurchases = user.Balance.SumPurchases
		userSpecs.Balance.SumBonuses = user.Balance.SumBonuses
	} else {
		userSpecs.Balance.SumPurchases = 0
		userSpecs.Balance.SumBonuses = 0
	}
	return &userSpecs
}
