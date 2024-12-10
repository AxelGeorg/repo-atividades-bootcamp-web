package repository

import (
	"app/internal"
)

type RepositoryTicketMock struct {
	FuncGet                            func() (map[int]internal.TicketAttributes, error)
	FuncGetTicketsByDestinationCountry func(country string) (map[int]internal.TicketAttributes, error)

	Spy struct {
		Get                            int
		GetTicketsByDestinationCountry int
	}
}

func NewRepositoryTicketMock() RepositoryTicketMock {
	return RepositoryTicketMock{}
}

func (r *RepositoryTicketMock) GetTotalTickets() (map[int]internal.TicketAttributes, error) {
	r.Spy.Get++
	return r.FuncGet()
}

func (r *RepositoryTicketMock) GetTicketsByDestinationCountry(country string) (map[int]internal.TicketAttributes, error) {
	r.Spy.GetTicketsByDestinationCountry++
	return r.FuncGetTicketsByDestinationCountry(country)
}
