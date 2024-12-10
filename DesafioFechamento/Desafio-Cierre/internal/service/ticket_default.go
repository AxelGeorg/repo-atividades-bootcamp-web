package service

import "app/internal/repository"

// ServiceTicketDefault represents the default service of the tickets
type ServiceTicketDefault struct {
	// rp represents the repository of the tickets
	Repository repository.RepositoryTicket
}

// NewServiceTicketDefault creates a new default service of the tickets
func NewServiceTicketDefault(rp repository.RepositoryTicket) ServiceTicketDefault {
	return ServiceTicketDefault{
		Repository: rp,
	}
}

// GetTotalTickets returns the total number of tickets
func (s *ServiceTicketDefault) GetTotalTickets() (total int, err error) {
	tickets, err := s.Repository.GetAll() // Obtém todos os tickets do repositório
	if err != nil {
		return 0, err // Retorna erro se a recuperação falhar
	}

	total = len(tickets) // Conta o número de tickets
	return total, nil    // Retorna o total e nil para o erro
}
