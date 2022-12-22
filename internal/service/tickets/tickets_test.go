package tickets

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	ticketsDomain "homework/internal/domain/tickets"
	mockTicketsService "homework/internal/service/tickets/mock"
	"homework/internal/util/terr"
)

//go:generate mockgen -destination ./mock/tickets_service_mock.go homework/internal/service/tickets TicketsService

func Test_CreateTicket(t *testing.T) {

	// Arrange
	ticketId := uuid.MustParse("6382589b-ab8e-4519-8c00-d0fe095179b3")
	flightId := uuid.MustParse("7d5925a6-2016-4c72-9298-517fc40d936c")
	userId := uuid.MustParse("07d87607-1f06-4599-8af5-07229525c106")
	passengerId := uuid.MustParse("b8d0b64d-08d8-4f9d-8c5c-cabd44957f16")
	classSeatId := uuid.MustParse("07d87607-1f06-4599-8af5-07229525c106")
	seatId := uuid.MustParse("c6eff2bf-525d-4b81-b995-d812874bbba8")
	paramsCreateTicket := &ticketsDomain.ParamsCreateTicket{
		StatusTimestamp:        time.Now(),
		FlightId:               flightId,
		UserId:                 userId,
		PassengerId:            &passengerId,
		ParamsCreatePassenger:  nil,
		ClassSeatsId:           classSeatId,
		SeatId:                 &seatId,
		CountAdditionalBaggage: 1,
		Price:                  1000,
	}

	var tests = []struct {
		name string
		args *ticketsDomain.ParamsCreateTicket
		want uuid.UUID
		err  error
	}{
		{
			name: "success",
			args: paramsCreateTicket,
			want: ticketId,
			err:  nil,
		},
		{
			name: "fail/sql database error",
			args: paramsCreateTicket,
			want: uuid.UUID{},
			err:  terr.SQLDatabaseError(errors.New("")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ticketsService := mockTicketsService.NewMockTicketsService(ctrl)
			ticketsService.EXPECT().
				CreateTicket(ctx, tt.args).
				Return(tt.want, tt.err)

			// Act
			got, err := ticketsService.CreateTicket(ctx, tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_PayForTicket(t *testing.T) {

	// Arrange
	ticketId := uuid.MustParse("6382589b-ab8e-4519-8c00-d0fe095179b3")
	paramsPayForTicket := &ticketsDomain.ParamsPayForTicket{
		TicketId: ticketId,
		UserId:   uuid.MustParse("07d87607-1f06-4599-8af5-07229525c106"),
	}

	var tests = []struct {
		name string
		args *ticketsDomain.ParamsPayForTicket
		want uuid.UUID
		err  error
	}{
		{
			name: "success",
			args: paramsPayForTicket,
			want: ticketId,
			err:  nil,
		},
		{
			name: "fail/sql database error",
			args: paramsPayForTicket,
			want: uuid.UUID{},
			err:  terr.SQLDatabaseError(errors.New("")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ticketsService := mockTicketsService.NewMockTicketsService(ctrl)
			ticketsService.EXPECT().
				PayForTicket(ctx, tt.args).
				Return(tt.want, tt.err)

			// Act
			got, err := ticketsService.PayForTicket(ctx, tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_RegisterTicket(t *testing.T) {

	// Arrange
	ticketId := uuid.MustParse("6382589b-ab8e-4519-8c00-d0fe095179b3")
	seatId := uuid.MustParse("c6eff2bf-525d-4b81-b995-d812874bbba8")
	paramsRegisterTicket := &ticketsDomain.ParamsRegisterTicket{
		TicketId: ticketId,
		SeatId:   &seatId,
	}

	var tests = []struct {
		name string
		args *ticketsDomain.ParamsRegisterTicket
		want uuid.UUID
		err  error
	}{
		{
			name: "success",
			args: paramsRegisterTicket,
			want: ticketId,
			err:  nil,
		},
		{
			name: "fail/sql database error",
			args: paramsRegisterTicket,
			want: uuid.UUID{},
			err:  terr.SQLDatabaseError(errors.New("")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ticketsService := mockTicketsService.NewMockTicketsService(ctrl)
			ticketsService.EXPECT().
				RegisterTicket(ctx, tt.args).
				Return(tt.want, tt.err)

			// Act
			got, err := ticketsService.RegisterTicket(ctx, tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
