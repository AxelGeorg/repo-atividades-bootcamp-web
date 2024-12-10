package service_test

import (
	"app/internal"
	"app/internal/repository"
	"app/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceTicketDefault_GetTicketsAmountByDestinationCountry(t *testing.T) {
	t.Run("success to get total tickets", func(t *testing.T) {
		country := "USA"

		rp := repository.NewRepositoryTicketMock()

		rp.FuncGetTicketsByDestinationCountry = func(country string) (map[int]internal.TicketAttributes, error) {
			if country == "USA" {
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
			}
			return nil, nil
		}

		sv := service.NewServiceTicketDefault(&rp)
		total, err := sv.GetTicketsAmountByDestinationCountry(country)

		expectedTotal := 2
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)

		require.Equal(t, 1, rp.Spy.GetTicketsByDestinationCountry)
	})
}
