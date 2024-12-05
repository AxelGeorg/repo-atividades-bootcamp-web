package storage

type Product struct {
	Id           string
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

/*
type Products interface {
	Get() (p []Product, err error)
	GetByID(id int) (p *Product, err error)
	Save(p *Product) (err error)
	UpdateOrCreate(p *Product) (err error)
}
*/
