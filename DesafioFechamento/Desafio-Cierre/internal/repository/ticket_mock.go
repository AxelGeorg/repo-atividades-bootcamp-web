package repository

import (
	"app/internal"
	"context"
)

func NewRepositoryTicketMock() *RepositoryTicketMap {
	return &RepositoryTicketMap{}
}

type RepositoryTicketMock struct {
	FuncGet                            func() (t map[int]internal.TicketAttributes, err error)
	FuncGetTicketsByDestinationCountry func(country string) (t map[int]internal.TicketAttributes, err error)

	Spy struct {
		Get                            int
		GetTicketsByDestinationCountry int
	}
}

func (r *RepositoryTicketMock) Get(ctx context.Context) (t map[int]internal.TicketAttributes, err error) {
	r.Spy.Get++
	t, err = r.FuncGet()
	return
}

func (r *RepositoryTicketMock) GetTicketsByDestinationCountry(ctx context.Context, country string) (t map[int]internal.TicketAttributes, err error) {
	r.Spy.GetTicketsByDestinationCountry++
	t, err = r.FuncGetTicketsByDestinationCountry(country)
	return
}
