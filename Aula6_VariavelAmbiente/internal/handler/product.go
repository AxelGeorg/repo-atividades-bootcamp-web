package handler

import (
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

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
	Data    *Data  `json:"data,omitempty"`
	Error   bool   `json:"error"`
}

type ResponseBodyTotalPrice struct {
	Products   []*Data `json:"products,omitempty"`
	TotalPrice float64 `json:"total_price"`
}

type Data struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ProductController struct {
	Service service.Service
}

func NewHandlerProducts(service service.Service) *ProductController {
	return &ProductController{
		Service: service,
	}
}

func (c *ProductController) handleError(w http.ResponseWriter, message string, statusCode int) {
	body := &ResponseBodyProduct{
		Message: message,
		Data:    nil,
		Error:   true,
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ProductController) respondWithProduct(w http.ResponseWriter, message string, product storage.Product) {
	dt := Data{
		Id:           product.Id,
		Name:         product.Name,
		Code_value:   product.Code_value,
		Is_published: product.Is_published,
		Expiration:   product.Expiration,
		Quantity:     product.Quantity,
		Price:        product.Price,
	}

	body := &ResponseBodyProduct{
		Message: message,
		Data:    &dt,
		Error:   false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ProductController) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"message": "` + message + `"}`))
}

func (c *ProductController) validateToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Token")
	if token == "" {
		c.respondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
		return false
	}

	if token != os.Getenv("TOKEN") {
		c.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	return true
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	var reqBody RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		c.handleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product := storage.Product{
		Name:         reqBody.Name,
		Quantity:     reqBody.Quantity,
		Code_value:   reqBody.Code_value,
		Is_published: reqBody.Is_published,
		Expiration:   reqBody.Expiration,
		Price:        reqBody.Price,
	}

	productServ, err := c.Service.Create(product)
	if err != nil {
		code := http.StatusBadRequest
		body := &ResponseBodyProduct{
			Message: "Bad Request",
			Data:    nil,
			Error:   true,
		}

		log.Println(err)

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
		return
	}

	dt := Data{
		Id:           productServ.Id,
		Name:         productServ.Name,
		Code_value:   productServ.Code_value,
		Is_published: productServ.Is_published,
		Expiration:   productServ.Expiration,
		Quantity:     productServ.Quantity,
		Price:        productServ.Price,
	}

	body := &ResponseBodyProduct{
		Message: "Product created",
		Data:    &dt,
		Error:   false,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ProductController) UpdateOrCreate(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	var reqBody RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		c.handleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	product := storage.Product{
		Id:           idStr,
		Name:         reqBody.Name,
		Quantity:     reqBody.Quantity,
		Code_value:   reqBody.Code_value,
		Is_published: reqBody.Is_published,
		Expiration:   reqBody.Expiration,
		Price:        reqBody.Price,
	}

	productServ, err := c.Service.Update(product)
	if err != nil {
		if err.Error() == "product not found" {
			productServ, err = c.Service.Create(product)
			if err != nil {
				c.handleError(w, "Could not create product", http.StatusInternalServerError)
				return
			}
			c.respondWithProduct(w, "Product created", productServ)
			return
		}

		c.handleError(w, "Could not update product", http.StatusBadRequest)
		return
	}

	c.respondWithProduct(w, "Product updated", productServ)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	idStr := chi.URLParam(r, "id")

	_, err := c.Service.GetById(idStr)
	if err != nil {
		c.handleError(w, "Product not found", http.StatusNotFound)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := c.Service.Patch(idStr, updates)
	if err != nil {
		switch err.Error() {
		case "product not found":
			http.Error(w, "Product not found", http.StatusNotFound)
		default:
			http.Error(w, "Could not update product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	idStr := chi.URLParam(r, "id")

	if _, err := c.Service.GetById(idStr); err != nil {
		c.handleError(w, "Product not found", http.StatusNotFound)
		return
	}

	err := c.Service.Delete(idStr)
	if err != nil {
		c.handleError(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	body := &ResponseBodyProduct{
		Message: "Product deleted",
		Data:    nil,
		Error:   false,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	products, err := c.Service.GetAll()
	if err != nil {
		c.handleError(w, "Could not retrieve products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	idStr := chi.URLParam(r, "id")

	product, err := c.Service.GetById(idStr)
	if err != nil {
		if err.Error() == "product not found" {
			c.handleError(w, "Product not found", http.StatusNotFound)
		} else {
			c.handleError(w, "Could not retrieve product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (c *ProductController) Search(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.handleError(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	products, err := c.Service.SearchByPrice(price)
	if err != nil {
		c.handleError(w, "Could not retrieve products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) ConsumerPrice(w http.ResponseWriter, r *http.Request) {
	if !c.validateToken(w, r) {
		return
	}

	listIds := r.URL.Query().Get("list")

	var ids []string
	if listIds != "" {
		ids = strings.Split(listIds, ",")
	}

	totalPrice, products, err := c.Service.GetTotalPrice(ids)
	if err != nil {
		c.handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var productsResponse []*Data
	for _, product := range products {
		dt := Data{
			Id:           product.Id,
			Name:         product.Name,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Quantity:     product.Quantity,
			Price:        product.Price,
		}

		productsResponse = append(productsResponse, &dt)
	}

	body := &ResponseBodyTotalPrice{
		Products:   productsResponse,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
