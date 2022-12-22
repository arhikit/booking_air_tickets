package tickets

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	flightsDomain "homework/internal/domain/flights"
	usersDomain "homework/internal/domain/users"
	"time"

	ticketsDomain "homework/internal/domain/tickets"
	"homework/internal/util/terr"
)

type TicketsStorage interface {
	GetPassengerById(ctx context.Context, passengerId uuid.UUID) (*ticketsDomain.Passenger, error)
	GetTicketById(ctx context.Context, ticketId uuid.UUID) (*ticketsDomain.Ticket, error)
	CreateTicket(ctx context.Context, paramsCreateTicket *ticketsDomain.ParamsCreateTicket) (uuid.UUID, error)
	PayForTicket(ctx context.Context, paramsPayForTicket *ticketsDomain.ParamsPayForTicket) (uuid.UUID, error)
	RefundTicket(ctx context.Context, paramsRefundTicket *ticketsDomain.ParamsRefundTicket) (uuid.UUID, error)
	RegisterTicket(ctx context.Context, paramsRegisterTicket *ticketsDomain.ParamsRegisterTicket) (uuid.UUID, error)
}

type storage struct {
	db *pgxpool.Pool
}

func (s storage) GetPassengerById(ctx context.Context, passengerId uuid.UUID) (*ticketsDomain.Passenger, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`
     		SELECT 	passenger.id, 
     		       
					users.id,
					users.name,
					users.email,
     		       	
					passenger.name_passenger,
     		       	passenger.identity_data_passenger
        		FROM passengers passenger

       			INNER JOIN users
     				ON passenger.user_id = users.id

			WHERE passenger.id = $1`,
		passengerId.String())

	var user usersDomain.User
	var passenger ticketsDomain.Passenger
	err = row.Scan(
		&passenger.Id,

		&user.Id,
		&user.Name,
		&user.Email,

		&passenger.NamePassenger,
		&passenger.IdentityDataPassenger,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found passenger (id %s)", passengerId))
		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}

	passenger.User = user
	return &passenger, nil
}

func (s storage) GetTicketById(ctx context.Context, ticketId uuid.UUID) (*ticketsDomain.Ticket, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`
     		SELECT 	ticket.id, 
						
					status.id,
                   	status.name,
                   	ticket.status_timestamp,
     		       
     			 	flight.id, 
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
     		        flight.pet_allowed,

					users.id,
					users.name,
					users.email,

     		 		passenger.id, 
     		       	passenger.name_passenger,
     		       	passenger.identity_data_passenger,

     				class_seats.id,
   	 				class_seats.name,
    				class_seats.count_seats,
    				class_seats.width,
    				class_seats.pitch,
    				class_seats.count_in_row,
    		       	
					CASE
						WHEN ticket.seat_id IS NOT NULL
							THEN true
						ELSE false
					END is_seat_assigned, 
					CASE
						WHEN ticket.seat_id IS NOT NULL
							THEN seat.id
						ELSE ticket.id
					END seat_id,
     		       	CASE
						WHEN ticket.seat_id IS NOT NULL
							THEN seat.number
						ELSE status.name
					END seat_number,

					ticket.count_additional_baggage,
					ticket.price,
					ticket.paid_with_bonuses,
					ticket.accrued_bonuses

       		FROM tickets ticket

      			INNER JOIN statuses status
     				ON ticket.status_id = status.id

      			INNER JOIN flights flight
     				ON ticket.flight_id = flight.id
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

      			INNER JOIN users
     				ON ticket.user_id = users.id

      			INNER JOIN passengers passenger
     				ON ticket.passenger_id = passenger.id

      			INNER JOIN classes_seats class_seats
     				ON ticket.class_seats_id = class_seats.id

      			LEFT JOIN seats seat
     				ON ticket.seat_id = seat.id

 			WHERE ticket.id = $1`,
		ticketId.String())

	var status ticketsDomain.Status
	var aircraft flightsDomain.Aircraft
	var airline flightsDomain.Airline
	var airportDeparture flightsDomain.Airport
	var cityDeparture flightsDomain.City
	var airportArrival flightsDomain.Airport
	var cityArrival flightsDomain.City
	var durationMin int
	var flight flightsDomain.Flight
	var user usersDomain.User
	var passenger ticketsDomain.Passenger
	var classSeats flightsDomain.ClassSeats
	var isSeatAssigned bool
	var seat flightsDomain.Seat
	var ticket ticketsDomain.Ticket

	err = row.Scan(
		&ticket.Id,

		&status.Id,
		&status.Name,
		&status.Timestamp,

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

		&user.Id,
		&user.Name,
		&user.Email,

		&passenger.Id,
		&passenger.NamePassenger,
		&passenger.IdentityDataPassenger,

		&classSeats.Id,
		&classSeats.Name,
		&classSeats.CountSeats,
		&classSeats.Width,
		&classSeats.Pitch,
		&classSeats.CountInRow,

		&isSeatAssigned,
		&seat.Id,
		&seat.Number,

		&ticket.CountAdditionalBaggage,
		&ticket.Price,
		&ticket.PaidWithBonuses,
		&ticket.AccruedBonuses,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, terr.NotFound(fmt.Sprintf("not found ticket (id %s)", ticketId))

		} else {
			return nil, terr.SQLDatabaseError(err)
		}
	}

	ticket.Status = status

	aircraft.Airline = airline
	airportDeparture.City = cityDeparture
	airportArrival.City = cityArrival
	flight.Aircraft = aircraft
	flight.DepartureAirport = airportDeparture
	flight.ArrivalAirport = airportArrival
	flight.Duration = time.Duration(durationMin) * time.Minute
	ticket.Flight = flight

	ticket.User = user

	passenger.User = user
	ticket.Passenger = passenger

	classSeats.Aircraft = aircraft
	ticket.ClassSeats = classSeats

	if isSeatAssigned {
		seat.ClassSeats = classSeats
		ticket.Seat = &seat
	}

	return &ticket, nil
}

func (s storage) CreateTicket(ctx context.Context, paramsCreateTicket *ticketsDomain.ParamsCreateTicket) (uuid.UUID, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	// в зависимости от параметров входящего post-запроса генерируем различные sql-запросы
	var sqlQuery string

	// указываем параметры запроса
	ticketId := uuid.New()
	arrParams := []interface{}{
		ticketId.String(),
		paramsCreateTicket.StatusTimestamp,
		paramsCreateTicket.FlightId.String(),
		paramsCreateTicket.UserId.String(),
		paramsCreateTicket.ClassSeatsId.String(),
		paramsCreateTicket.CountAdditionalBaggage,
		paramsCreateTicket.Price,
	}

	// указываем текст запроса для создания билета.
	if paramsCreateTicket.SeatId != nil {
		// создание билета с выбранным местом
		sqlQuery = `
	 		INSERT INTO tickets (
	 		            	id,
							status_id,
							status_timestamp,
	 		                flight_id,
	 		                user_id,
	 		                class_seats_id,
	 		                count_additional_baggage,
	 		                price,
	 		                paid_with_bonuses,
	 		                accrued_bonuses,
							seat_id,
	 		            	passenger_id
	 				)
	 				VALUES (
	 						$1,
							1,
	 				        $2,
	 				        $3,
	 				        $4,
	 				        $5,
	 				        $6,
	 				        $7,
							0,
							0,
							$8,
	 				        $9
	 				)`
	} else {
		// создание билета без выбранного места
		sqlQuery = `
	 		INSERT INTO tickets (
	 		            	id,
							status_id,
							status_timestamp,
	 		                flight_id,
	 		                user_id,
	 		                class_seats_id,
	 		                count_additional_baggage,
	 		                price,
	 		                paid_with_bonuses,
	 		                accrued_bonuses,
	 		            	passenger_id
	 				)
	 				VALUES (
	 						$1,
							1,
	 				        $2,
	 				        $3,
	 				        $4,
	 				        $5,
	 				        $6,
	 				        $7,
							0,
							0,
							$8
					)`
	}

	// добавляем параметры запроса в зависимости от того,
	// будет или нет создаваться пассажир и будет ли указываться место.
	// при создании пассажира добавляется запрос на создание пассажира
	if paramsCreateTicket.PassengerId == nil && paramsCreateTicket.SeatId != nil {
		// создание пассажира и создание билета с выбранным местом
		arrParams = append(arrParams,
			paramsCreateTicket.SeatId.String(),
			uuid.New().String(),
			paramsCreateTicket.ParamsCreatePassenger.NamePassenger,
			paramsCreateTicket.ParamsCreatePassenger.IdentityDataPassenger,
		)
		sqlQuery = `WITH A AS (INSERT INTO passengers (
	 		            	id,
	 		                user_id,
	 		                name_passenger,
	 		                identity_data_passenger
	 					)
	 					VALUES (
	 						$9,
	 				        $4,
	 				        $10,
	 				        $11
	 					))` + sqlQuery
	} else if paramsCreateTicket.PassengerId == nil && paramsCreateTicket.SeatId == nil {
		// создание пассажира и создание билета без выбранного места
		arrParams = append(arrParams,
			uuid.New().String(),
			paramsCreateTicket.ParamsCreatePassenger.NamePassenger,
			paramsCreateTicket.ParamsCreatePassenger.IdentityDataPassenger,
		)
		sqlQuery = `WITH A AS (INSERT INTO passengers (
	 		            	id,
	 		                user_id,
	 		                name_passenger,
	 		                identity_data_passenger
	 					)
	 					VALUES (
	 						$8,
	 				        $4,
	 				        $9,
	 				        $10
	 					))` + sqlQuery
	} else if paramsCreateTicket.PassengerId != nil && paramsCreateTicket.SeatId != nil {
		// используется существующий пассажир и создание билета с выбранным местом
		arrParams = append(arrParams,
			paramsCreateTicket.SeatId.String(),
			paramsCreateTicket.PassengerId.String(),
		)
	} else if paramsCreateTicket.PassengerId != nil && paramsCreateTicket.SeatId == nil {
		// используется существующий пассажир и создание билета без выбранного места
		arrParams = append(arrParams,
			paramsCreateTicket.PassengerId.String(),
		)
	}

	rows, err := conn.Query(ctx, sqlQuery, arrParams...)
	defer rows.Close()

	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}
	return ticketId, nil
}

func (s storage) PayForTicket(ctx context.Context, paramsPayForTicket *ticketsDomain.ParamsPayForTicket) (uuid.UUID, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	// в зависимости от параметров входящего post-запроса генерируем различные sql-запросы
	var sqlQuery string

	// Билету (tickets) устанавливаются:
	// - статус status_id = 2(Paid) и время изменения статуса status_timestamp
	// - сумма начисляемых бонусных баллов accrued_bonuses
	// - сумма бонусов, использованных для оплаты билета paid_with_bonuses
	sqlQuery = `
     		UPDATE tickets
				SET status_id = 2, 
					status_timestamp = $2, 
					paid_with_bonuses = $3, 
					accrued_bonuses = $4 
   			WHERE id = $1`

	// указываем параметры запроса
	arrParams := []interface{}{
		paramsPayForTicket.TicketId.String(),
		paramsPayForTicket.StatusTimestamp,
		paramsPayForTicket.PaidWithBonuses,
		paramsPayForTicket.AccruedBonuses,
	}

	// Если для пользователя еще не заполнен баланс, то добавляется запись в таблицу users_balance.
	// Сумма покупок sum_purchases устанавливается равной стоимости билета.
	if !paramsPayForTicket.UserBalanceInit {
		arrParams = append(arrParams,
			uuid.New().String(),
			paramsPayForTicket.UserId.String(),
			paramsPayForTicket.Price,
		)
		sqlQuery = `WITH A AS (INSERT INTO users_balance (
		            	id,
		                user_id,
		                sum_purchases,
		                sum_bonuses
					)
					VALUES (
						$5,
				        $6,
				        $7,
				        0
					))` + sqlQuery

	} else {
		// Изменения баланса пользователя (users_balance):
		// - по пользователю увеличивается общая сумма покупок sum_purchases на стоимость билета.
		// - по пользователю уменьшается общая сумма бонусов sum_bonuses на сумму бонусов, использованную при покупке билета.
		arrParams = append(arrParams,
			paramsPayForTicket.UserId.String(),
			paramsPayForTicket.Price,
			paramsPayForTicket.PaidWithBonuses,
		)
		sqlQuery = `WITH A AS (UPDATE users_balance
					SET sum_purchases = sum_purchases + $6, 
						sum_bonuses = sum_bonuses - $7 
					WHERE user_id = $5)` + sqlQuery
	}

	rows, err := conn.Query(ctx, sqlQuery, arrParams...)
	defer rows.Close()

	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}

	ticketId := paramsPayForTicket.TicketId
	return ticketId, nil
}

func (s storage) RefundTicket(ctx context.Context, paramsRefundTicket *ticketsDomain.ParamsRefundTicket) (uuid.UUID, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	// генерируем sql-запрос
	var sqlQuery string

	// Билету (tickets) устанавливаются:
	// - статус status_id = 4(Refunded) и время изменения статуса status_timestamp
	sqlQuery = `
     		UPDATE tickets
				SET status_id = 4, 
					status_timestamp = $2
   			WHERE id = $1`

	// указываем параметры запроса
	arrParams := []interface{}{
		paramsRefundTicket.TicketId.String(),
		paramsRefundTicket.StatusTimestamp,
	}

	// Изменения баланса пользователя (users_balance):
	// - по пользователю уменьшается общая сумма покупок на стоимость билета.
	// - по пользователю увеличивается общая сумма бонусов на стоимость билета.
	// Таким образом, возвращаются на баланс пользователя
	// и сумма бонусов, использованная при покупке билета paid_with_bonuses,
	// и сумма оплаченных денег за билет ticket.Price-PaidWithBonuses.
	arrParams = append(arrParams,
		paramsRefundTicket.UserId.String(),
		paramsRefundTicket.Price,
	)
	sqlQuery = `WITH A AS (UPDATE users_balance
					SET sum_purchases = sum_purchases - $4, 
						sum_bonuses = sum_bonuses + $4 
					WHERE user_id = $3)` + sqlQuery

	rows, err := conn.Query(ctx, sqlQuery, arrParams...)
	defer rows.Close()

	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}

	ticketId := paramsRefundTicket.TicketId
	return ticketId, nil
}

func (s storage) RegisterTicket(ctx context.Context, paramsRegisterTicket *ticketsDomain.ParamsRegisterTicket) (uuid.UUID, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}
	defer conn.Release()

	// в зависимости от параметров входящего post-запроса генерируем различные sql-запросы
	var sqlQuery string

	// указываем параметры запроса
	arrParams := []interface{}{
		paramsRegisterTicket.UserId.String(),
		paramsRegisterTicket.AccruedBonuses,
		paramsRegisterTicket.TicketId.String(),
		paramsRegisterTicket.StatusTimestamp,
	}

	// Изменения баланса пользователя (users_balance):
	// - по пользователю увеличивается общая сумма бонусов на сумму начисленных за билет бонусов accrued_bonuses
	sqlQuery = `WITH A AS (UPDATE users_balance
					SET sum_bonuses = sum_bonuses + $2 
					WHERE user_id = $1)`

	// Билету (tickets) устанавливаются:
	// - статус status_id = 5(Registered) и время изменения статуса status_timestamp
	// - место, если при покупке билета место не было назначено
	if paramsRegisterTicket.SeatId != nil {
		arrParams = append(arrParams,
			paramsRegisterTicket.SeatId,
		)
		sqlQuery = sqlQuery + `
     		UPDATE tickets
				SET status_id = 5, 
					status_timestamp = $4,
				    seat_id = $5
   			WHERE id = $3`
	} else {
		sqlQuery = sqlQuery + `
     		UPDATE tickets
				SET status_id = 5, 
					status_timestamp = $4
   			WHERE id = $3`
	}

	rows, err := conn.Query(ctx, sqlQuery, arrParams...)
	defer rows.Close()

	if err != nil {
		return uuid.UUID{}, terr.SQLDatabaseError(err)
	}

	ticketId := paramsRegisterTicket.TicketId
	return ticketId, nil
}

func NewTicketsStorage(db *pgxpool.Pool) TicketsStorage {
	return &storage{db: db}
}