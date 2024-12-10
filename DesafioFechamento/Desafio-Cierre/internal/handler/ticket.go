package handler

type TicketHandler struct {
	Service service.ti
}

func NewHandlerProducts(service service.Service) *ProductController {
	return &ProductController{
		Service: service,
	}
}
