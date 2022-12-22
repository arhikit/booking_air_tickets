package v1

import (
	"encoding/json"
	"homework/internal/util/terr"
	"net/http"

	"homework/specs"
)

func (a apiServer) GetFlights(w http.ResponseWriter, r *http.Request, paramsGetFlightsSpecs specs.GetFlightsParams) {

	paramsGetFlights, err := transformParamsGetFlights(&paramsGetFlightsSpecs)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ctx := r.Context()
	flights, err := a.serviceRegistry.Flight.GetFlights(ctx, paramsGetFlights)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	flightsSpecs := make([]specs.Flight, len(flights))
	for i, flight := range flights {
		flightsSpecs[i] = *transformFlight(&flight)
	}
	_ = json.NewEncoder(w).Encode(flightsSpecs)
}

func (a apiServer) GetFlightById(w http.ResponseWriter, r *http.Request, flightIdSpecs specs.UUIDPathObjectID) {

	flightId, err := convertStringToUuid(string(flightIdSpecs))
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_FLIGHT_UUID", err.Error()))
		return
	}

	ctx := r.Context()
	flight, err := a.serviceRegistry.Flight.GetFlightById(ctx, flightId)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	flightSpecs := transformFlight(flight)
	_ = json.NewEncoder(w).Encode(flightSpecs)
}

func (a apiServer) GetFlightVacantSeats(w http.ResponseWriter, r *http.Request, flightIdSpecs specs.UUIDPathObjectID) {
	flightId, err := convertStringToUuid(string(flightIdSpecs))
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_FLIGHT_UUID", err.Error()))
		return
	}

	ctx := r.Context()
	arrVacantSeats, err := a.serviceRegistry.Flight.GetFlightVacantSeats(ctx, flightId)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	arrVacantSeatsSpecs := make([]specs.VacantSeats, len(arrVacantSeats))
	for i, vacantSeats := range arrVacantSeats {
		arrVacantSeatsSpecs[i] = *transformVacantSeats(&vacantSeats)
	}
	_ = json.NewEncoder(w).Encode(arrVacantSeatsSpecs)
}
