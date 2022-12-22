package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"homework/internal/util/terr"
	"net/http"

	specs "homework/specs"
)

func (a apiServer) GetTicketById(w http.ResponseWriter, r *http.Request, ticketIdSpecs specs.UUIDPathObjectID) {

	ticketId, err := convertStringToUuid(string(ticketIdSpecs))
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_TICKET_UUID", err.Error()))
		return
	}

	ctx := r.Context()
	ticket, err := a.serviceRegistry.Ticket.GetTicketById(ctx, ticketId)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ticketSpecs := transformTicket(ticket)
	_ = json.NewEncoder(w).Encode(ticketSpecs)

}

func (a apiServer) CreateTicket(w http.ResponseWriter, r *http.Request) {

	paramsCreateTicketSpecs := &specs.ParamsCreateTicket{}
	err := json.NewDecoder(r.Body).Decode(paramsCreateTicketSpecs)
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_BODY_REQUEST", err.Error()))
		return
	}

	paramsCreateTicket, err := transformParamsCreateTicket(paramsCreateTicketSpecs)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ctx := r.Context()
	ticketId, err := a.serviceRegistry.Ticket.CreateTicket(ctx, paramsCreateTicket)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	createdItem := specs.CreatedItem{Id: uuid.UUID(ticketId).String()}
	_ = json.NewEncoder(w).Encode(createdItem)

}

func (a apiServer) PayForTicket(w http.ResponseWriter, r *http.Request) {

	paramsPayForTicketSpecs := &specs.ParamsPayForTicket{}
	err := json.NewDecoder(r.Body).Decode(paramsPayForTicketSpecs)
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_BODY_REQUEST", err.Error()))
		return
	}

	paramsPayForTicket, err := transformParamsPayForTicket(paramsPayForTicketSpecs)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ctx := r.Context()
	ticketId, err := a.serviceRegistry.Ticket.PayForTicket(ctx, paramsPayForTicket)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	updatedItem := specs.UpdatedItem{Id: uuid.UUID(ticketId).String()}
	_ = json.NewEncoder(w).Encode(updatedItem)

}

func (a apiServer) RefundTicket(w http.ResponseWriter, r *http.Request) {

	paramsRefundTicketSpecs := &specs.ParamsRefundTicket{}
	err := json.NewDecoder(r.Body).Decode(paramsRefundTicketSpecs)
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_BODY_REQUEST", err.Error()))
		return
	}

	paramsRefundTicket, err := transformParamsRefundTicket(paramsRefundTicketSpecs)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ctx := r.Context()
	ticketId, err := a.serviceRegistry.Ticket.RefundTicket(ctx, paramsRefundTicket)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	updatedItem := specs.UpdatedItem{Id: uuid.UUID(ticketId).String()}
	_ = json.NewEncoder(w).Encode(updatedItem)
}

func (a apiServer) RegisterTicket(w http.ResponseWriter, r *http.Request) {

	paramsRegisterTicketSpecs := &specs.ParamsRegisterTicket{}
	err := json.NewDecoder(r.Body).Decode(paramsRegisterTicketSpecs)
	if err != nil {
		terr.WriteError(w, terr.BadRequest("INVALID_BODY_REQUEST", err.Error()))
		return
	}

	paramsRegisterTicket, err := transformParamsRegisterTicket(paramsRegisterTicketSpecs)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	ctx := r.Context()
	ticketId, err := a.serviceRegistry.Ticket.RegisterTicket(ctx, paramsRegisterTicket)
	if err != nil {
		terr.WriteError(w, err.(*terr.Error))
		return
	}

	updatedItem := specs.UpdatedItem{Id: uuid.UUID(ticketId).String()}
	_ = json.NewEncoder(w).Encode(updatedItem)

}
