package repository

import (
	"app/internal"
	"context"
)

// RepositoryTicket represents the repository interface for tickets
type RepositoryTicket interface {
	// GetAll returns all the tickets
	GetAll() (map[int]internal.TicketAttributes, error)

	// GetTicketByDestinationCountry returns the tickets filtered by destination country
	GetTicketsByDestinationCountry(ctx context.Context, country string) (map[int]internal.TicketAttributes, error)
}
