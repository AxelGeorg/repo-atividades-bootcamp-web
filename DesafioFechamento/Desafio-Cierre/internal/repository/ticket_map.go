package repository

import (
	"app/internal"
)

func NewRepositoryTicketMap(storage map[int]*internal.TicketAttributes) RepositoryTicketMap {
	return RepositoryTicketMap{
		db:     storage,
		lastId: len(storage) - 1,
	}
}

type RepositoryTicketMap struct {
	db     map[int]*internal.TicketAttributes
	lastId int
}

func (r *RepositoryTicketMap) GetTotalTickets() (map[int]internal.TicketAttributes, error) {
	t := make(map[int]internal.TicketAttributes, len(r.db))
	for k, v := range r.db {
		t[k] = *v
	}

	return t, nil
}

func (r *RepositoryTicketMap) GetTicketsByDestinationCountry(country string) (map[int]internal.TicketAttributes, error) {
	t := make(map[int]internal.TicketAttributes)
	for k, v := range r.db {
		if v.Country == country {
			t[k] = *v
		}
	}

	return t, nil
}
