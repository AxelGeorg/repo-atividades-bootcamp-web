package service_test

import (
	"app/internal"
	"app/internal/repository"
	"app/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceTicketDefault_GetTicketsAmountByDestinationCountry(t *testing.T) {
	tests := []struct {
		name          string
		country       string
		expectedTotal int
	}{
		{
			name:          "USA with 2 tickets",
			country:       "USA",
			expectedTotal: 2,
		},
		{
			name:          "Canada with 1 ticket",
			country:       "Canada",
			expectedTotal: 1,
		},
		{
			name:          "Mexico with 0 tickets",
			country:       "Mexico",
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := repository.NewRepositoryTicketMock()
			rp.FuncGetTicketsByDestinationCountry = func(country string) (map[int]internal.TicketAttributes, error) {
				switch country {
				case "USA":
					return map[int]internal.TicketAttributes{
						1: {
							Name:    "John",
							Email:   "johndoe@gmail.com",
							Country: "USA",
							Hour:    "10:00",
							Price:   100,
						},
						2: {
							Name:    "Jane",
							Email:   "janedoe@gmail.com",
							Country: "USA",
							Hour:    "11:00",
							Price:   150,
						},
					}, nil
				case "Canada":
					return map[int]internal.TicketAttributes{
						1: {
							Name:    "Alex",
							Email:   "alex@email.com",
							Country: "Canada",
							Hour:    "09:00",
							Price:   120,
						},
					}, nil
				default:
					return nil, nil
				}
			}

			sv := service.NewServiceTicketDefault(&rp)
			total, err := sv.GetTicketsAmountByDestinationCountry(tt.country)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTotal, total)

			require.Equal(t, 1, rp.Spy.GetTicketsByDestinationCountry)
		})
	}
}
