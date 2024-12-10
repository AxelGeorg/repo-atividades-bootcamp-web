package repository

import (
	"app/internal"
	"context"
)

// NewRepositoryTicketMap creates a new repository for tickets in a map
func NewRepositoryTicketMap(storage map[int]*internal.TicketAttributes) RepositoryTicketMap {
	return RepositoryTicketMap{
		db:     storage,
		lastId: len(storage) - 1,
	}
}

// RepositoryTicketMap implements the repository interface for tickets in a map
type RepositoryTicketMap struct {
	// db represents the database in a map
	// - key: id of the ticket
	// - value: ticket
	db map[int]*internal.TicketAttributes

	// lastId represents the last id of the ticket
	lastId int
}

// GetAll returns all the tickets
func (r *RepositoryTicketMap) GetAll() (map[int]internal.TicketAttributes, error) {
	// create a copy of the map
	t := make(map[int]internal.TicketAttributes, len(r.db))
	for k, v := range r.db {
		t[k] = *v
	}

	return t, nil
}

// GetTicketsByDestinationCountry returns the tickets filtered by destination country
func (r *RepositoryTicketMap) GetTicketsByDestinationCountry(ctx context.Context, country string) (map[int]internal.TicketAttributes, error) {
	// create a copy of the map
	t := make(map[int]internal.TicketAttributes)
	for k, v := range r.db {
		if v.Country == country {
			t[k] = *v
		}
	}

	return t, nil
}
