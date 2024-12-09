package handler

import (
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"aula4/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published *bool   `json:"is_published"`
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

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	var reqBody RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	if reqBody.Is_published == nil {
		falseValue := false
		reqBody.Is_published = &falseValue
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
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	RespondWithProduct(w, &productServ, http.StatusCreated, utils.MessageProductCreated)
}

func (c *ProductController) UpdateOrCreate(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	var reqBody RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

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
				ResponseWithError(w, err, http.StatusInternalServerError)
				return
			}

			RespondWithProduct(w, &productServ, http.StatusCreated, utils.MessageProductCreated)
			return
		}

		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	RespondWithProduct(w, &productServ, http.StatusOK, utils.MessageProductUpdated)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

	_, err = c.Service.GetById(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusNotFound)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	product, err := c.Service.Patch(idStr, updates)
	if err != nil {
		switch err.Error() {
		case "product not found":
			ResponseWithError(w, err, http.StatusNotFound)
		default:
			ResponseWithError(w, err, http.StatusInternalServerError)
		}
		return
	}

	RespondWithProduct(w, product, http.StatusOK, utils.MessageProductUpdated)
}

func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

	if _, err := c.Service.GetById(idStr); err != nil {
		ResponseWithError(w, err, http.StatusNotFound)
		return
	}

	err = c.Service.Delete(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithProduct(w, nil, http.StatusNoContent, utils.MessageProductDeleted)
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	products, err := c.Service.GetAll()
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

	var product *storage.Product
	product, err = c.Service.GetById(idStr)
	if err != nil {
		if err.Error() == "product not found" {
			ResponseWithError(w, err, http.StatusNotFound)
		} else {
			ResponseWithError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (c *ProductController) Search(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	products, err := c.Service.SearchByPrice(price)
	if err != nil {
		ResponseWithError(w, errors.New("could not retrieve products"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) ConsumerPrice(w http.ResponseWriter, r *http.Request) {
	if !ValidateToken(w, r) {
		return
	}

	listIds := r.URL.Query().Get("list")

	var ids []string
	if listIds != "" {
		ids = strings.Split(listIds, ",")
	}

	totalPrice, products, err := c.Service.GetTotalPrice(ids)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	var productsResponse []*Data
	for _, product := range products {
		dt := Data{
			Id:           product.Id,
			Name:         product.Name,
			Code_value:   product.Code_value,
			Is_published: *product.Is_published,
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
