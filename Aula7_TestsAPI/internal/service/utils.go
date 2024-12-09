package service

import (
	"aula4/internal/repository/storage"
	"aula4/internal/utils"
	"errors"
	"log"
	"strings"
)

func GetProductQuantity(service Service, ids []string) (int, []*storage.Product, error) {
	mapProdQtd := make(map[string]int)

	var quantity int
	var products []*storage.Product
	for _, id := range ids {
		idStr := strings.TrimSpace(id)
		err := utils.ValidateUUID(idStr)
		if err != nil {
			return 0.0, nil, err
		}

		product, err := service.GetById(idStr)
		if err != nil {
			log.Printf("Error retrieving product with ID:" + idStr + " - Error: " + err.Error())
			continue
		}

		mapProdQtd[idStr]++
		if mapProdQtd[idStr] > product.Quantity {
			return 0.0, nil, errors.New("not enough stock for product ID:" + idStr)
		}

		*product.Is_published = true
		_, err = service.Update(*product)
		if err != nil {
			return 0.0, nil, err
		}

		quantity++
		products = append(products, product)
	}

	if len(products) == 0 {
		return 0.0, nil, errors.New("no products available")
	}

	return quantity, products, nil
}

func GetProductQuantityTotal(service Service) (int, []*storage.Product, error) {
	products, err := service.GetAll()
	if err != nil {
		return 0.0, nil, err
	}

	var quantity int
	for _, product := range products {
		quantity += product.Quantity
	}

	return quantity, products, nil
}
