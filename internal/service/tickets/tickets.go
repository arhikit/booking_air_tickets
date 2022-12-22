package tickets

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	flightsDomain "homework/internal/domain/flights"
	ticketsDomain "homework/internal/domain/tickets"
	usersDomain "homework/internal/domain/users"
	"homework/internal/util/terr"
)

type TicketsService interface {
	GetTicketById(ctx context.Context, ticketId uuid.UUID) (*ticketsDomain.Ticket, error)
	CreateTicket(ctx context.Context, paramsCreateTicket *ticketsDomain.ParamsCreateTicket) (uuid.UUID, error)
	PayForTicket(ctx context.Context, paramsPayForTicket *ticketsDomain.ParamsPayForTicket) (uuid.UUID, error)
	RefundTicket(ctx context.Context, paramsRefundTicket *ticketsDomain.ParamsRefundTicket) (uuid.UUID, error)
	RegisterTicket(ctx context.Context, paramsRegisterTicket *ticketsDomain.ParamsRegisterTicket) (uuid.UUID, error)
}

type TicketsStorage interface {
	GetPassengerById(ctx context.Context, passengerId uuid.UUID) (*ticketsDomain.Passenger, error)
	GetTicketById(ctx context.Context, ticketId uuid.UUID) (*ticketsDomain.Ticket, error)
	CreateTicket(ctx context.Context, paramsCreateTicket *ticketsDomain.ParamsCreateTicket) (uuid.UUID, error)
	PayForTicket(ctx context.Context, paramsPayForTicket *ticketsDomain.ParamsPayForTicket) (uuid.UUID, error)
	RefundTicket(ctx context.Context, paramsRefundTicket *ticketsDomain.ParamsRefundTicket) (uuid.UUID, error)
	RegisterTicket(ctx context.Context, paramsRegisterTicket *ticketsDomain.ParamsRegisterTicket) (uuid.UUID, error)
}

type FlightsStorage interface {
	GetFlightById(ctx context.Context, flightId uuid.UUID) (*flightsDomain.Flight, error)
	GetFlightVacantSeatsByClassId(ctx context.Context, flightId uuid.UUID, classSeatsId uuid.UUID) (*flightsDomain.VacantSeats, error)
}

type UsersStorage interface {
	GetAccruedBonuses(ctx context.Context, userId uuid.UUID, ticketPrice int) (int, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error)
}

type service struct {
	ticketsStorage TicketsStorage
	flightsStorage FlightsStorage
	usersStorage   UsersStorage
}

func (s service) GetTicketById(ctx context.Context, ticketId uuid.UUID) (*ticketsDomain.Ticket, error) {
	return s.ticketsStorage.GetTicketById(ctx, ticketId)
}

func (s service) CreateTicket(ctx context.Context, paramsCreateTicket *ticketsDomain.ParamsCreateTicket) (uuid.UUID, error) {

	// проверяем, что по переданному FlightId существует рейс
	flight, err := s.flightsStorage.GetFlightById(ctx, paramsCreateTicket.FlightId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки рейса:
	// до вылета осталось больше 2 часов
	if flight.DepartureDate.Sub(paramsCreateTicket.StatusTimestamp).Hours() < 2 {
		return uuid.UUID{}, terr.BadRequest("FLIGHT_ALREADY_CLOSED", "sale of tickets for the flight is closed")
	}

	// проверяем, что по переданному UserId существует пользователь
	_, err = s.usersStorage.GetUserById(ctx, paramsCreateTicket.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// если пассажир уже существует, то проверяем его
	if paramsCreateTicket.PassengerId != nil {

		// проверяем, что по переданному PassengerId существует пассажир
		passenger, err := s.ticketsStorage.GetPassengerById(ctx, *paramsCreateTicket.PassengerId)
		if err != nil {
			return uuid.UUID{}, err
		}

		// проверки пассажира:
		// пользователь пассажира соответствует пользователю билета
		if paramsCreateTicket.UserId != passenger.User.Id {
			return uuid.UUID{}, terr.BadRequest("INVALID_PASSENGER", fmt.Sprintf("the passenger's user (id %s) doesn't match the user of the ticket (id %s)", passenger.User.Id, paramsCreateTicket.UserId))
		}

		// проверки того, что на данного пассажира уже может быть билет на данный рейс нет.
		// причина: в ржд можно купить на одного пассажира несколько билетов, чтобы выкупить полностью купе.
		// предполагаю, что и на самолет можно купить несколько билетов на одного пассажира, чтобы выкупить весь ряд или весь самолет.
	}

	// проверяем, что на данном рейсе существуют места с заданным классом ClassSeatsId
	vacantSeats, err := s.flightsStorage.GetFlightVacantSeatsByClassId(ctx, paramsCreateTicket.FlightId, paramsCreateTicket.ClassSeatsId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки класса места:
	// есть свободные места данного класса
	if vacantSeats.CountVacantSeats == 0 {
		return uuid.UUID{}, terr.BadRequest("NO_VACANT_SEAT", fmt.Sprintf("no vacant seats with class seat (id %s) ", paramsCreateTicket.ClassSeatsId))
	}

	// если место было указано, то проверяем его
	if paramsCreateTicket.SeatId != nil {

		// проверяем, что место есть в списке свободных мест
		isSeatVacant := false
		seatId := *paramsCreateTicket.SeatId
		for _, seat := range vacantSeats.Seats {
			if seat.Id == seatId {
				isSeatVacant = true
				break
			}
		}

		// место занято
		if !isSeatVacant {
			return uuid.UUID{}, terr.BadRequest("SEAT_DOESNT_VACANT", fmt.Sprintf("seat (id %s) isn't in the list of vacant seats", seatId))
		}
	}

	// рассчитаем стоимость билета как сумму стоимости билета выбранного класса
	// + стоимость дополнительного багажа * количество дополнительного багажа
	// + стоимость выбора места, если место было выбрано на этапе создания билета
	var price int
	for _, flightPrice := range flight.PricesTickets {
		if flightPrice.ClassSeats.Id == paramsCreateTicket.ClassSeatsId {
			price = flightPrice.PriceTicket
			break
		}
	}
	price += paramsCreateTicket.CountAdditionalBaggage * flight.PriceAdditionalBaggage
	if paramsCreateTicket.SeatId != nil {
		price += flight.PriceSeatSelection
	}
	paramsCreateTicket.Price = price

	// создаем билет и пассажира, если он не существует
	ticketId, err := s.ticketsStorage.CreateTicket(ctx, paramsCreateTicket)
	return ticketId, err
}

func (s service) PayForTicket(ctx context.Context, paramsPayForTicket *ticketsDomain.ParamsPayForTicket) (uuid.UUID, error) {

	// по id получаем билет для оплаты
	ticket, err := s.ticketsStorage.GetTicketById(ctx, paramsPayForTicket.TicketId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки билета:
	// оплатить можно только новый билет со статусом 1 (Created)
	if ticket.Status.Id != 1 {
		return uuid.UUID{}, terr.BadRequest("INVALID_STATUS_TICKET", fmt.Sprintf("ticket (id %s) has wrong status (%s)", paramsPayForTicket.TicketId, ticket.Status.Name))
	}

	// оплатить можно только билет, созданный менее 15 мин назад, иначе билет должен быть отменен
	if paramsPayForTicket.StatusTimestamp.Sub(ticket.Status.Timestamp).Minutes() > 15 {
		return uuid.UUID{}, terr.BadRequest("TICKET_ALREADY_CANCELED", "time to pay is over")
	}

	// проверяем, что по переданному UserId существует пользователь
	user, err := s.usersStorage.GetUserById(ctx, paramsPayForTicket.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки пользователя:
	// переданный пользователь соответствует пользователю билета
	if paramsPayForTicket.UserId != ticket.User.Id {
		return uuid.UUID{}, terr.BadRequest("INVALID_USER", fmt.Sprintf("the user (id %s) doesn't match the user of the ticket (id %s)", paramsPayForTicket.UserId, ticket.User.Id))
	}

	// проверки, если передается сумма бонусов для оплаты
	if paramsPayForTicket.PaidWithBonuses > 0 {

		// проверяем, что у пользователя достаточно бонусов
		if user.Balance == nil || user.Balance.SumBonuses < paramsPayForTicket.PaidWithBonuses {
			return uuid.UUID{}, terr.BadRequest("INVALID_SUM_BONUSES", "user doesn't have enough bonuses")
		}

		// проверяем, что переданная сумма бонусов не превышает половину стоимости билета
		if paramsPayForTicket.PaidWithBonuses > int(ticket.Price/2) {
			return uuid.UUID{}, terr.BadRequest("INVALID_SUM_BONUSES", "sum bonuses is more than half of the ticket price")
		}
	}

	// Все проверки пройдены

	// Здесь по логике бизнес-процесса выполняется обращение к платежной системе
	// и производится оплата на сумму ticket.Price-PaidWithBonuses

	// Получаем сумму бонусных баллов AccruedBonuses, начисляемых за приобретение билета.
	// Бонусные баллы поступят на счет пользователя только после регистрации на рейс. До этого момента информация о них хранится только в билете.
	// Расчет бонусов - % от суммы общей покупок пользователя.
	accruedBonuses, err := s.usersStorage.GetAccruedBonuses(ctx, paramsPayForTicket.UserId, ticket.Price)
	if err != nil {
		return uuid.UUID{}, err
	}
	paramsPayForTicket.AccruedBonuses = accruedBonuses

	// передаем стоимость билета для изменения баланса пользователя
	paramsPayForTicket.Price = ticket.Price
	paramsPayForTicket.UserBalanceInit = user.Balance != nil

	// Выполняем изменение билета, в т.ч. начисление бонусов за билет, и изменение баланса пользователя
	ticketId, err := s.ticketsStorage.PayForTicket(ctx, paramsPayForTicket)
	return ticketId, err
}

func (s service) RefundTicket(ctx context.Context, paramsRefundTicket *ticketsDomain.ParamsRefundTicket) (uuid.UUID, error) {

	// по id получаем билет для возврата
	ticket, err := s.ticketsStorage.GetTicketById(ctx, paramsRefundTicket.TicketId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки билета:
	// вернуть можно только оплаченный билет со статусом 2 (Paid)
	if ticket.Status.Id != 2 {
		return uuid.UUID{}, terr.BadRequest("INVALID_STATUS_TICKET", fmt.Sprintf("ticket (id %s) has wrong status (%s)", paramsRefundTicket.TicketId, ticket.Status.Name))
	}

	// вернуть билет можно только в случае, если до вылета осталось больше 24 часов
	if ticket.Flight.DepartureDate.Sub(paramsRefundTicket.StatusTimestamp).Hours() < 24 {
		return uuid.UUID{}, terr.BadRequest("REFUND_ALREADY_CLOSED", "flight ticket refund is not possible")
	}

	// проверяем, что по переданному UserId существует пользователь
	user, err := s.usersStorage.GetUserById(ctx, paramsRefundTicket.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки пользователя:
	// переданный пользователь соответствует пользователю билета
	if paramsRefundTicket.UserId != ticket.User.Id {
		return uuid.UUID{}, terr.BadRequest("INVALID_USER", fmt.Sprintf("the user (id %s) doesn't match the user of the ticket (id %s)", paramsRefundTicket.UserId, ticket.User.Id))
	}

	// баланс пользователя должен быть заполнен, т.к. данный билет уже был куплен и это должно быть отражено в балансе пользователя
	if user.Balance == nil {
		return uuid.UUID{}, terr.BadRequest("INVALID_USER", "no information about the user's balance")
	}

	// Все проверки пройдены

	// передаем стоимость билета для изменения баланса пользователя
	paramsRefundTicket.Price = ticket.Price

	// Выполняем изменение билета и изменение баланса пользователя
	ticketId, err := s.ticketsStorage.RefundTicket(ctx, paramsRefundTicket)
	return ticketId, err
}

func (s service) RegisterTicket(ctx context.Context, paramsRegisterTicket *ticketsDomain.ParamsRegisterTicket) (uuid.UUID, error) {

	// по id получаем билет для выполнения регистрации
	ticket, err := s.ticketsStorage.GetTicketById(ctx, paramsRegisterTicket.TicketId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки билета:
	// зарегистрировать можно только оплаченный билет со статусом 2 (Paid)
	if ticket.Status.Id != 2 {
		return uuid.UUID{}, terr.BadRequest("INVALID_STATUS_TICKET", fmt.Sprintf("ticket (id %s) has wrong status (%s)", paramsRegisterTicket.TicketId, ticket.Status.Name))
	}

	// зарегистрировать билет можно только в случае, если до вылета осталось больше 1 часа и меньше 24 часов
	if ticket.Flight.DepartureDate.Sub(paramsRegisterTicket.StatusTimestamp).Hours() > 24 {
		return uuid.UUID{}, terr.BadRequest("CHECK_IN_DOESNT_START", "check-in hasn't started yet")
	} else if ticket.Flight.DepartureDate.Sub(paramsRegisterTicket.StatusTimestamp).Hours() < 1 {
		return uuid.UUID{}, terr.BadRequest("CHECK_IN_ALREADY_CLOSED", "check-in is already closed")
	}

	// проверяем, что по переданному UserId существует пользователь
	user, err := s.usersStorage.GetUserById(ctx, paramsRegisterTicket.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// проверки пользователя:
	// переданный пользователь соответствует пользователю билета
	if paramsRegisterTicket.UserId != ticket.User.Id {
		return uuid.UUID{}, terr.BadRequest("INVALID_USER", fmt.Sprintf("the user (id %s) doesn't match the user of the ticket (id %s)", paramsRegisterTicket.UserId, ticket.User.Id))
	}

	// баланс пользователя должен быть заполнен, т.к. данный билет уже был куплен и это должно быть отражено в балансе пользователя
	if user.Balance == nil {
		return uuid.UUID{}, terr.BadRequest("INVALID_USER", "no information about the user's balance")
	}

	// проверяем место, если при покупке билета место не было назначено
	if ticket.Seat == nil {

		// проверяем, что передано место.
		// т.к. в билете место еще не заполнено, значит, место должно назначаться при регистрации на рейс
		if paramsRegisterTicket.SeatId == nil {
			return uuid.UUID{}, terr.BadRequest("SEAT_DOESNT_ASSIGNED", "seat has to be assigned")
		}

		// проверяем, что на данном рейсе существуют место SeatId
		// получаем список вакантных места рейса с заданным классом ClassSeatsId
		vacantSeats, err := s.flightsStorage.GetFlightVacantSeatsByClassId(ctx, ticket.Flight.Id, ticket.ClassSeats.Id)
		if err != nil {
			return uuid.UUID{}, err
		}

		// проверяем, что место есть в списке свободных мест
		isSeatVacant := false
		seatId := *paramsRegisterTicket.SeatId
		for _, seat := range vacantSeats.Seats {
			if seat.Id == seatId {
				isSeatVacant = true
				break
			}
		}
		// место занято
		if !isSeatVacant {
			return uuid.UUID{}, terr.BadRequest("SEAT_DOESNT_VACANT", fmt.Sprintf("seat (id %s) isn't in the list of vacant seats", seatId))
		}
	}

	// Все проверки пройдены

	// передаем сумму начисленных за билет бонусов для изменения баланса пользователя
	paramsRegisterTicket.AccruedBonuses = ticket.AccruedBonuses

	// Выполняем изменение билета и изменение баланса пользователя
	ticketId, err := s.ticketsStorage.RegisterTicket(ctx, paramsRegisterTicket)
	return ticketId, err
}

func NewTicketsService(ticketsStorage TicketsStorage, flightsStorage FlightsStorage, usersStorage UsersStorage) TicketsService {
	return &service{
		ticketsStorage: ticketsStorage,
		flightsStorage: flightsStorage,
		usersStorage:   usersStorage,
	}
}
