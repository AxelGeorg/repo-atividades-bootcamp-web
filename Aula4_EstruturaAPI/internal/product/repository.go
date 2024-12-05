package product

import (
	"aula4/internal/domain"
	"encoding/json"
	"fmt"
	"os"
)

func FillProductList(db *map[int]*domain.Product) {
	file, err := os.Open("products.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	var products []domain.Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, product := range products {
		(*db)[product.Id] = &product
	}
}
