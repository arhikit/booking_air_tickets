package flights

import (
	"context"
	"github.com/google/uuid"

	flightsDomain "homework/internal/domain/flights"
)

type service struct {
	flightsStorage FlightsStorage
}

type FlightsService interface {
	GetFlights(ctx context.Context, paramsGetFlights *flightsDomain.ParamsGetFlights) ([]flightsDomain.Flight, error)
	GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error)
	GetFlightVacantSeats(ctx context.Context, flightId uuid.UUID) ([]flightsDomain.VacantSeats, error)
}

type FlightsStorage interface {
	GetCityById(ctx context.Context, cityId uuid.UUID) (*flightsDomain.City, error)
	GetFlights(ctx context.Context, paramsGetFlights *flightsDomain.ParamsGetFlights) ([]flightsDomain.Flight, error)
	GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error)
	GetFlightVacantSeats(ctx context.Context, flightId uuid.UUID) ([]flightsDomain.VacantSeats, error)
}

func (s service) GetFlights(ctx context.Context, paramsGetFlights *flightsDomain.ParamsGetFlights) ([]flightsDomain.Flight, error) {

	// проверяем, что по переданному DepartureCityId существует город
	_, err := s.flightsStorage.GetCityById(ctx, paramsGetFlights.DepartureCityId)
	if err != nil {
		return nil, err
	}

	// проверяем, что по переданному ArrivalCityId существует город
	_, err = s.flightsStorage.GetCityById(ctx, paramsGetFlights.ArrivalCityId)
	if err != nil {
		return nil, err
	}

	return s.flightsStorage.GetFlights(ctx, paramsGetFlights)
}

func (s service) GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error) {
	return s.flightsStorage.GetFlightById(ctx, flightId)
}

func (s service) GetFlightVacantSeats(ctx context.Context, flightId uuid.UUID) ([]flightsDomain.VacantSeats, error) {
	return s.flightsStorage.GetFlightVacantSeats(ctx, flightId)
}

func NewFlightsService(flightsStorage FlightsStorage) FlightsService {
	return &service{flightsStorage: flightsStorage}
}
