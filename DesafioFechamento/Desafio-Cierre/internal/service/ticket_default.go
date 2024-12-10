package service

import "app/internal/repository"

type ServiceTicketDefault struct {
	Repository repository.RepositoryTicket
}

func NewServiceTicketDefault(rp repository.RepositoryTicket) ServiceTicketDefault {
	return ServiceTicketDefault{
		Repository: rp,
	}
}

func (s *ServiceTicketDefault) GetTotalTickets() (int, error) {
	tickets, err := s.Repository.GetTotalTickets()
	if err != nil {
		return 0, err
	}

	total := len(tickets)
	return total, nil
}

func (s *ServiceTicketDefault) GetTicketsAmountByDestinationCountry(destination string) (int, error) {
	tickets, err := s.Repository.GetTicketsByDestinationCountry(destination)
	if err != nil {
		return 0, err
	}

	total := len(tickets)
	return total, nil
}

func (s *ServiceTicketDefault) GetPercentageTicketsByDestinationCountry(destination string) (float64, error) {
	ticketsByDest, err := s.Repository.GetTicketsByDestinationCountry(destination)
	if err != nil {
		return 0, err
	}

	ticketsTotal, err := s.GetTotalTickets()
	if err != nil {
		return 0, err
	}

	percentage := (float64(len(ticketsByDest)) / float64(ticketsTotal)) * 100
	return percentage, nil
}
