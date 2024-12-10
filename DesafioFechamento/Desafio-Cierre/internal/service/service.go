package service

type ServiceTicket interface {
	GetTotalTickets() (int, error)

	GetTicketsAmountByDestinationCountry(destination string) (int, error)

	GetPercentageTicketsByDestinationCountry(destination string) (float64, error)
}
