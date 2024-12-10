package service

type ServiceTicket interface {
	// GetTotalAmountTickets returns the total amount of tickets
	GetTotalAmountTickets() (total int, err error)

	// GetTicketsAmountByDestinationCountry returns the amount of tickets filtered by destination country
	// ...

	// GetPercentageTicketsByDestinationCountry returns the percentage of tickets filtered by destination country
	// ...
}
