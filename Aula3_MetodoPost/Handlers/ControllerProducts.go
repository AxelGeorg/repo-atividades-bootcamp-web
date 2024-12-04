package handlers

import "net/http"

type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

type ControllerProducts struct {
	storage map[int]*Product
}

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string `json:"message"`
	Data    *struct {
		Id           int     `json:"id"`
		Name         string  `json:"name"`
		Quantity     int     `json:"quantity"`
		Code_value   string  `json:"code_value"`
		Is_published bool    `json:"is_published"`
		Expiration   string  `json:"expiration"`
		Price        float64 `json:"price"`
	} `json:"data"`
	Error bool `json:"error"`
}

func (c *ControllerProducts) Create() http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
