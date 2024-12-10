package repository

import (
	"app/internal"
)

type RepositoryTicket interface {
	GetTotalTickets() (map[int]internal.TicketAttributes, error)

	GetTicketsByDestinationCountry(country string) (map[int]internal.TicketAttributes, error)
}
