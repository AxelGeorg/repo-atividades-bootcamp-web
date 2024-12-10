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
		country := "USA" // Define o país que estamos testando

		// Criando o mock para o repositório de bilhetes
		rp := repository.NewRepositoryTicketMock() // Certifique-se de usar o mock correto

		// Configurando o FuncGetTicketsByDestinationCountry para retornar bilhetes para "USA"
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

		// Criando o serviço utilizando o repositório mock
		sv := service.NewServiceTicketDefault(&rp)

		// Chamando a função que estamos testando
		total, err := sv.GetTicketsAmountByDestinationCountry(country)

		// O valor esperado é 2, já que temos dois bilhetes para o país "USA"
		expectedTotal := 2
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)

		// Verifica se o método foi chamado
		require.Equal(t, 1, rp.Spy.GetTicketsByDestinationCountry)
	})
}
