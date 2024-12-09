package storage

type Storage interface {
	ReadAllProductsToFile() ([]*Product, error)
	WriteProductsToFile(productList []*Product) error

	ReadProductById(id string) (*Product, error)
	SaveProduct(product *Product) error
	UpdateProduct(updatedProduct *Product) error
	DeleteProduct(id string) error
}
