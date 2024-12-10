package handler

import (
	"app/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseBodyTicket struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func ResponseWithError(w http.ResponseWriter, err error, statusCode int) {
	body := &ResponseBodyTicket{
		Message: http.StatusText(statusCode) + " - " + err.Error(),
		Error:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func RespondWithSuccess(w http.ResponseWriter, statusCode int, message string) {
	body := ResponseBodyTicket{
		Message: message,
		Error:   false,
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
	country := r.URL.Path[len("/ticket/getByCountry/"):]
	fmt.Println(country)

	countTickets, err := c.Service.GetTicketsAmountByDestinationCountry(country)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithSuccess(w, http.StatusOK, fmt.Sprint("How many people are traveling to ", country, " is ", countTickets))
}

func (c *TicketHandler) GetAverage(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Path[len("/ticket/getAverage/"):]
	fmt.Println(country)

	percentage, err := c.Service.GetPercentageTicketsByDestinationCountry(country)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithSuccess(w, http.StatusOK, fmt.Sprint("Percentage of people traveling to ", country, " is ", percentage, "%"))
}
