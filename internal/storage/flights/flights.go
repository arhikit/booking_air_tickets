package flights

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"homework/internal/util/terr"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	flightsDomain "homework/internal/domain/flights"
)

type FlightsStorage interface {
	GetCityById(ctx context.Context, cityId uuid.UUID) (*flightsDomain.City, error)
	GetFlights(ctx context.Context, paramsGetFlights *flightsDomain.ParamsGetFlights) ([]flightsDomain.Flight, error)
	GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error)
	GetFlightVacantSeats(ctx context.Context, flightId uuid.UUID) ([]flightsDomain.VacantSeats, error)
	GetFlightVacantSeatsByClassId(ctx context.Context, flightId uuid.UUID, classSeatsId uuid.UUID) (*flightsDomain.VacantSeats, error)
}

type storage struct {
	db *pgxpool.Pool
}

// получение данных рейсов

func getSqlQueryFlights(sqlQueryCondition string) string {
	return `SELECT 	flight.id, 
     		       	flight.name,
					
					aircraft.id,
					aircraft.name,     		         		       
     		       		airline.id,
     		       		airline.name,
					
					airport_departure.id,
					airport_departure.name,
     		       		city_departure.id,
     		       		city_departure.name,
     		       	
     		       	airport_arrival.id,
     		       	airport_arrival.name,    		          		       
     		       		city_arrival.id,
     		       		city_arrival.name,  
     		        
					flight.departure_date,
     		        flight.duration,
     		       
					flight.price_additional_baggage,
     		        flight.price_seat_selection,
     		        flight.is_international,
     		        flight.baggage_included,
     		        flight.pet_allowed
     		FROM flights flight
      			INNER JOIN aircrafts aircraft
     				ON flight.aircraft_id = aircraft.id
     		    	INNER JOIN airlines airline
     					ON aircraft.airline_id = airline.id
     		    
    			INNER JOIN airports airport_departure
     				ON flight.departure_airport_id = airport_departure.id
     				INNER JOIN cities city_departure
     					ON airport_departure.city_id = city_departure.id

    			INNER JOIN airports airport_arrival
     				ON flight.arrival_airport_id = airport_arrival.id
     				INNER JOIN cities city_arrival
     					ON airport_arrival.city_id = city_arrival.id
			WHERE ` + sqlQueryCondition
}

func scanFlight(row pgx.Row) (flightsDomain.Flight, error) {

	var airline flightsDomain.Airline
	var aircraft flightsDomain.Aircraft
	var cityDeparture flightsDomain.City
	var airportDeparture flightsDomain.Airport
	var cityArrival flightsDomain.City
	var airportArrival flightsDomain.Airport
	var durationMin int
	var flight flightsDomain.Flight

	err := row.Scan(
		&flight.Id,
		&flight.Name,

		&aircraft.Id,
		&aircraft.Name,
		&airline.Id,
		&airline.Name,

		&airportDeparture.Id,
		&airportDeparture.Name,
		&cityDeparture.Id,
		&cityDeparture.Name,

		&airportArrival.Id,
		&airportArrival.Name,
		&cityArrival.Id,
		&cityArrival.Name,

		&flight.DepartureDate,
		&durationMin,

		&flight.PriceAdditionalBaggage,
		&flight.PriceSeatSelection,
		&flight.IsInternational,
		&flight.BaggageIncluded,
		&flight.PetAllowed,
	)

	if err != nil {
		return flight, err
	}
	aircraft.Airline = airline
	airportDeparture.City = cityDeparture
	airportArrival.City = cityArrival

	flight.Aircraft = aircraft
	flight.DepartureAirport = airportDeparture
	flight.ArrivalAirport = airportArrival

	flight.Duration = time.Duration(durationMin) * time.Minute
	return flight, nil
}

func (s storage) getFlightPrices(ctx context.Context, SqlQueryCondition string, paramsQuery []interface{}) (map[uuid.UUID][]flightsDomain.FlightPrice, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`WITH selected_flights AS (SELECT 
				flights_prices.flight_id flight_id,
				flights_prices.class_seats_id class_seats_id,   			
				flights_prices.price_ticket price_ticket
			FROM flights_prices
      			INNER JOIN flights flight
     				ON flights_prices.flight_id = flight.id
				INNER JOIN airports airport_departure
					ON flight.departure_airport_id = airport_departure.id
				INNER JOIN airports airport_arrival
					ON flight.arrival_airport_id = airport_arrival.id
			WHERE `+SqlQueryCondition+`) 
     		SELECT 	
    				selected_flights.flight_id,
    				class_seats.id,
					class_seats.name,
    				class_seats.count_seats,
    				class_seats.width,
    				class_seats.pitch,
    				class_seats.count_in_row,
					aircraft.id,
					aircraft.name,
   		       		airline.id,
   		       		airline.name,
    				selected_flights.price_ticket,
					class_seats.count_seats - CASE
							WHEN busy_class_seats.count_busy IS NOT NULL
								THEN busy_class_seats.count_busy
							ELSE 0
						END AS count_vacant
					
     		FROM selected_flights
        	    INNER JOIN classes_seats class_seats
       				ON selected_flights.class_seats_id = class_seats.id
					INNER JOIN aircrafts aircraft
						ON class_seats.aircraft_id = aircraft.id
						INNER JOIN airlines airline
							ON aircraft.airline_id = airline.id

        		LEFT JOIN (SELECT
        		                tickets.flight_id,
        		                tickets.class_seats_id,
        		                SUM(1) AS count_busy
        		            FROM tickets
        		            	INNER JOIN selected_flights
									ON selected_flights.flight_id = tickets.flight_id
       								AND selected_flights.class_seats_id = tickets.class_seats_id
          		            WHERE tickets.status_id <> 3 and tickets.status_id <> 4
       		            	GROUP BY
        		                tickets.flight_id,
        		                tickets.class_seats_id) busy_class_seats
 	   				ON selected_flights.flight_id = busy_class_seats.flight_id
	   					AND selected_flights.class_seats_id = busy_class_seats.class_seats_id`,
		paramsQuery...)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer rows.Close()

	mapFlightsPrices := make(map[uuid.UUID][]flightsDomain.FlightPrice)
	for rows.Next() {

		var flightId uuid.UUID
		var airline flightsDomain.Airline
		var aircraft flightsDomain.Aircraft
		var classSeats flightsDomain.ClassSeats
		var flightPrice flightsDomain.FlightPrice

		err = rows.Scan(
			&flightId,
			&classSeats.Id,
			&classSeats.Name,
			&classSeats.CountSeats,
			&classSeats.Width,
			&classSeats.Pitch,
			&classSeats.CountInRow,
			&aircraft.Id,
			&aircraft.Name,
			&airline.Id,
			&airline.Name,
			&flightPrice.PriceTicket,
			&flightPrice.CountVacantSeats,
		)

		if err != nil {
			return nil, terr.SQLDatabaseError(err)
		}

		aircraft.Airline = airline
		classSeats.Aircraft = aircraft
		flightPrice.ClassSeats = classSeats

		flightPrices := mapFlightsPrices[flightId]
		flightPrices = append(flightPrices, flightPrice)
		mapFlightsPrices[flightId] = flightPrices
	}
	return mapFlightsPrices, nil
}

func (s storage) GetCityById(ctx context.Context, cityId uuid.UUID) (*flightsDomain.City, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`
     		SELECT city.id, 
     		       city.name
        	FROM cities city 
 			WHERE city.id = $1`,
		cityId.String())

	var city flightsDomain.City
	err = row.Scan(
		&city.Id,
		&city.Name,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found city (id %s)", cityId))

		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}
	return &city, nil
}

func (s storage) GetFlights(ctx context.Context, paramsGetFlights *flightsDomain.ParamsGetFlights) ([]flightsDomain.Flight, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	paramsQuery := []interface{}{
		paramsGetFlights.DepartureCityId.String(),
		paramsGetFlights.ArrivalCityId.String(),
		paramsGetFlights.DepartureDate,
	}

	sqlQueryCondition := `airport_departure.city_id = $1 
							AND airport_arrival.city_id = $2
							AND flight.departure_date::date = $3`

	mapFlightsPrices, err := s.getFlightPrices(ctx, sqlQueryCondition, paramsQuery)
	if err != nil {
		return nil, err
	}

	sqlQuery := getSqlQueryFlights(sqlQueryCondition)
	rows, err := conn.Query(ctx, sqlQuery, paramsQuery...)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer rows.Close()

	var flights []flightsDomain.Flight
	for rows.Next() {

		flight, err := scanFlight(rows)
		if err != nil {
			return nil, terr.SQLDatabaseError(err)
		}
		flight.PricesTickets = mapFlightsPrices[flight.Id]

		flights = append(flights, flight)
	}
	return flights, nil
}

func (s storage) GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	paramsQuery := []interface{}{
		flightId.String(),
	}

	sqlQueryCondition := "flight.id = $1"
	sqlQuery := getSqlQueryFlights(sqlQueryCondition)
	row := conn.QueryRow(ctx, sqlQuery, paramsQuery...)

	flight, err := scanFlight(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found flight (id %s)", flightId))
		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}

	mapFlightsPrices, err := s.getFlightPrices(ctx, sqlQueryCondition, paramsQuery)
	if err != nil {
		return nil, err
	}
	flight.PricesTickets = mapFlightsPrices[flight.Id]

	return &flight, err
}

// получение свободных мест рейса в разрезе классов

func getSqlQueryVacantSeats(sqlQueryCondition string) string {
	return `WITH selected_classes_seats AS (SELECT 
				flight.id flight_id,
				class_seats.id class_seats_id,   			
				class_seats.name class_seats_name,   			
				class_seats.count_seats class_seats_count   			
        	FROM classes_seats class_seats      	    
        	    INNER JOIN flights flight
       				ON class_seats.aircraft_id = flight.aircraft_id
			WHERE ` + sqlQueryCondition + `)
			SELECT 
     		    selected_classes_seats.class_seats_id,
     		    selected_classes_seats.class_seats_name,
				selected_classes_seats.class_seats_count - CASE
						WHEN busy_class_seats.count_busy IS NOT NULL
							THEN busy_class_seats.count_busy
						ELSE 0
					END AS count_vacant
        	FROM selected_classes_seats
        		LEFT JOIN (SELECT
        		                ticket.flight_id,
        		                ticket.class_seats_id,
        		                SUM(1) AS count_busy
        		            FROM tickets ticket
        	    				INNER JOIN selected_classes_seats
       								ON ticket.flight_id = selected_classes_seats.flight_id
       								AND ticket.class_seats_id = selected_classes_seats.class_seats_id
         		            WHERE ticket.status_id <> 3 and ticket.status_id <> 4
       		            	GROUP BY
        		                ticket.flight_id,
        		                ticket.class_seats_id) busy_class_seats
 	   				ON selected_classes_seats.flight_id = busy_class_seats.flight_id
	   					AND selected_classes_seats.class_seats_id = busy_class_seats.class_seats_id`
}

func scanVacantSeats(row pgx.Row) (flightsDomain.VacantSeats, error) {

	var vacantSeats flightsDomain.VacantSeats
	err := row.Scan(
		&vacantSeats.ClassSeatsId,
		&vacantSeats.ClassSeatsName,
		&vacantSeats.CountVacantSeats,
	)
	if err != nil {
		return vacantSeats, err
	}
	return vacantSeats, nil
}

func (s storage) getMapVacantSeats(ctx context.Context, SqlQueryCondition string, paramsQuery []interface{}) (map[uuid.UUID][]flightsDomain.Seat, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT 	
					seat.id,
					seat.number,	
     				seat.class_seats_id
				FROM seats seat
					INNER JOIN classes_seats class_seats
						ON seat.class_seats_id = class_seats.id
					INNER JOIN flights flight
						ON class_seats.aircraft_id = flight.aircraft_id
					LEFT JOIN tickets ticket
						ON ticket.flight_id = flight.id
						AND ticket.seat_id = seat.id
			WHERE ticket.id IS NULL AND `+SqlQueryCondition,
		paramsQuery...)

	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer rows.Close()

	mapVacantSeats := make(map[uuid.UUID][]flightsDomain.Seat)
	for rows.Next() {

		var classSeatsId uuid.UUID
		var seat flightsDomain.Seat

		err = rows.Scan(
			&seat.Id,
			&seat.Number,
			&classSeatsId,
		)

		if err != nil {
			return nil, terr.SQLDatabaseError(err)
		}

		vacantSeats := mapVacantSeats[classSeatsId]
		vacantSeats = append(vacantSeats, seat)
		mapVacantSeats[classSeatsId] = vacantSeats
	}
	return mapVacantSeats, nil
}

func (s storage) GetFlightVacantSeats(ctx context.Context, flightId uuid.UUID) ([]flightsDomain.VacantSeats, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	paramsQuery := []interface{}{
		flightId.String(),
	}

	sqlQueryCondition := `flight.id = $1`
	mapVacantSeats, err := s.getMapVacantSeats(ctx, sqlQueryCondition, paramsQuery)
	if err != nil {
		return nil, err
	}

	sqlQuery := getSqlQueryVacantSeats(sqlQueryCondition)
	rows, err := conn.Query(ctx, sqlQuery, paramsQuery...)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer rows.Close()

	var arrVacantSeats []flightsDomain.VacantSeats
	for rows.Next() {

		vacantSeats, err := scanVacantSeats(rows)
		if err != nil {
			return nil, terr.SQLDatabaseError(err)
		}
		vacantSeats.Seats = mapVacantSeats[vacantSeats.ClassSeatsId]
		arrVacantSeats = append(arrVacantSeats, vacantSeats)
	}
	return arrVacantSeats, nil
}

func (s storage) GetFlightVacantSeatsByClassId(ctx context.Context, flightId uuid.UUID, classSeatsId uuid.UUID) (*flightsDomain.VacantSeats, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	paramsQuery := []interface{}{
		flightId.String(),
		classSeatsId.String(),
	}

	sqlQueryCondition := "flight.id = $1 and class_seats.id = $2"
	sqlQuery := getSqlQueryVacantSeats(sqlQueryCondition)
	row := conn.QueryRow(ctx, sqlQuery, paramsQuery...)

	vacantSeats, err := scanVacantSeats(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found class seat (id %s) in flight (id %s)", classSeatsId, flightId))
		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}

	mapVacantSeats, err := s.getMapVacantSeats(ctx, sqlQueryCondition, paramsQuery)
	if err != nil {
		return nil, err
	}
	vacantSeats.Seats = mapVacantSeats[vacantSeats.ClassSeatsId]

	return &vacantSeats, err
}

func NewFlightsStorage(db *pgxpool.Pool) FlightsStorage {
	return &storage{db: db}
}
