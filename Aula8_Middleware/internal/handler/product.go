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

type ProductController struct {
	Service service.Service
}

func NewHandlerProducts(service service.Service) *ProductController {
	return &ProductController{
		Service: service,
	}
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody utils.RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
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
		utils.ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	utils.RespondWithProduct(w, &productServ, http.StatusCreated, utils.MessageProductCreated)
}

func (c *ProductController) UpdateOrCreate(w http.ResponseWriter, r *http.Request) {
	var reqBody utils.RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
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
				utils.ResponseWithError(w, err, http.StatusInternalServerError)
				return
			}

			utils.RespondWithProduct(w, &productServ, http.StatusCreated, utils.MessageProductCreated)
			return
		}

		utils.ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	utils.RespondWithProduct(w, &productServ, http.StatusOK, utils.MessageProductUpdated)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
	}

	_, err = c.Service.GetById(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusNotFound)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	product, err := c.Service.Patch(idStr, updates)
	if err != nil {
		switch err.Error() {
		case "product not found":
			utils.ResponseWithError(w, err, http.StatusNotFound)
		default:
			utils.ResponseWithError(w, err, http.StatusInternalServerError)
		}
		return
	}

	utils.RespondWithProduct(w, product, http.StatusOK, utils.MessageProductUpdated)
}

func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
	}

	if _, err := c.Service.GetById(idStr); err != nil {
		utils.ResponseWithError(w, err, http.StatusNotFound)
		return
	}

	err = c.Service.Delete(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondWithProduct(w, nil, http.StatusNoContent, utils.MessageProductDeleted)
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := c.Service.GetAll()
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	err := utils.ValidateUUID(idStr)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
	}

	var product *storage.Product
	product, err = c.Service.GetById(idStr)
	if err != nil {
		if err.Error() == "product not found" {
			utils.ResponseWithError(w, err, http.StatusNotFound)
		} else {
			utils.ResponseWithError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (c *ProductController) Search(w http.ResponseWriter, r *http.Request) {
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		utils.ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	products, err := c.Service.SearchByPrice(price)
	if err != nil {
		utils.ResponseWithError(w, errors.New("could not retrieve products"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) ConsumerPrice(w http.ResponseWriter, r *http.Request) {
	listIds := r.URL.Query().Get("list")

	var ids []string
	if listIds != "" {
		ids = strings.Split(listIds, ",")
	}

	totalPrice, products, err := c.Service.GetTotalPrice(ids)
	if err != nil {
		utils.ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	var productsResponse []*utils.Data
	for _, product := range products {
		dt := utils.Data{
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

	body := &utils.ResponseBodyTotalPrice{
		Products:   productsResponse,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
