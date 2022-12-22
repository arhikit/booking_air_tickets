package v1

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	ticketsDomain "homework/internal/domain/tickets"
	"homework/internal/util/terr"
	"homework/specs"
)

func Test_ConvertStringToUuid(t *testing.T) {

	// Arrange
	validUuidString := "6ac8fa15-a3d7-4b5f-a6e5-5bce49da4647"
	validUuid := uuid.MustParse("6ac8fa15-a3d7-4b5f-a6e5-5bce49da4647")

	invalidUuidString := "123"
	invalidUuid, errInvalidUuid := uuid.Parse(invalidUuidString)

	var tests = []struct {
		name string
		args string
		want uuid.UUID
		err  error
	}{
		{
			name: "success",
			args: validUuidString,
			want: validUuid,
			err:  nil,
		},
		{
			name: "fail/empty id",
			args: "",
			want: uuid.UUID{},
			err:  errors.New("empty uuid"),
		},
		{
			name: "fail/invalid uuid",
			args: invalidUuidString,
			want: invalidUuid,
			err:  errInvalidUuid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Act
			got, err := convertStringToUuid(tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}

}

func Test_TransformParamsPayForTicket(t *testing.T) {

	// Arrange
	validParamsPayForTicketSpecs := &specs.ParamsPayForTicket{
		TicketId:        "6ac8fa15-a3d7-4b5f-a6e5-5bce49da4647",
		UserId:          "fdef87aa-7694-47c6-a5cd-50984326a071",
		PaidWithBonuses: 100,
	}
	validParamsPayForTicket := &ticketsDomain.ParamsPayForTicket{
		TicketId:        uuid.MustParse("6ac8fa15-a3d7-4b5f-a6e5-5bce49da4647"),
		UserId:          uuid.MustParse("fdef87aa-7694-47c6-a5cd-50984326a071"),
		PaidWithBonuses: 100,
	}

	invalidTicketIdString := "123"
	_, errInvalidTicketId := uuid.Parse(invalidTicketIdString)
	invalidParamsPayForTicketSpecs := &specs.ParamsPayForTicket{TicketId: invalidTicketIdString}

	var tests = []struct {
		name string
		args *specs.ParamsPayForTicket
		want *ticketsDomain.ParamsPayForTicket
		err  error
	}{
		{
			name: "success",
			args: validParamsPayForTicketSpecs,
			want: validParamsPayForTicket,
			err:  nil,
		},
		{
			name: "fail/invalid ParamsPayForTicket",
			args: invalidParamsPayForTicketSpecs,
			want: &ticketsDomain.ParamsPayForTicket{},
			err:  terr.BadRequest("INVALID_TICKET_UUID", errInvalidTicketId.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Act
			got, err := transformParamsPayForTicket(tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			tt.want.StatusTimestamp = got.StatusTimestamp
			assert.Equal(t, tt.want, got)
		})
	}

}
