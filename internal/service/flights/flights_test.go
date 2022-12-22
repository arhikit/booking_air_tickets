package flights

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	flightsDomain "homework/internal/domain/flights"
	mockFlightsService "homework/internal/service/flights/mock"
	"homework/internal/util/terr"
)

//go:generate mockgen -destination ./mock/flights_service_mock.go homework/internal/service/flights FlightsService

func Test_GetFlightById(t *testing.T) {

	// Arrange
	flightId := uuid.MustParse("7d5925a6-2016-4c72-9298-517fc40d936c")

	var tests = []struct {
		name string
		args uuid.UUID
		want *flightsDomain.Flight
		err  error
	}{
		{
			name: "success",
			args: flightId,
			want: &flightsDomain.Flight{
				Id: flightId,
			},
			err: nil,
		},
		{
			name: "fail/flight not found",
			args: flightId,
			want: nil,
			err:  terr.NotFound(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			flightsService := mockFlightsService.NewMockFlightsService(ctrl)
			flightsService.EXPECT().
				GetFlightById(ctx, tt.args).
				Return(tt.want, tt.err)

			// Act
			got, err := flightsService.GetFlightById(ctx, tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
