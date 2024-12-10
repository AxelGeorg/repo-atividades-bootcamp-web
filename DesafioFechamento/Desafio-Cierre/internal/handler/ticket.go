package handler

import (
	"app/internal"
	"app/internal/service"
	"encoding/json"
	"net/http"
)

type ResponseBodyTicket struct {
	Message string                     `json:"message"`
	Data    *internal.TicketAttributes `json:"data,omitempty"`
	Error   bool                       `json:"error"`
}

func ResponseWithError(w http.ResponseWriter, err error, statusCode int) {
	body := &ResponseBodyTicket{
		Message: http.StatusText(statusCode) + " - " + err.Error(),
		Data:    nil,
		Error:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

type TicketHandler struct {
	Service service.ServiceTicket
}

func NewHandlerTickets(service service.ServiceTicket) TicketHandler {
	return TicketHandler{
		Service: service,
	}
}

func (c *TicketHandler) GetByCountry(w http.ResponseWriter, r *http.Request) {
	tickets, err := c.Service.GetAll()
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}

func (c *TicketHandler) GetAverage(w http.ResponseWriter, r *http.Request) {
	tickets, err := c.Service.GetAll()
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}
