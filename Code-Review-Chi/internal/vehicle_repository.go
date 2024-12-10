package internal

// VehicleRepository is an interface that represents a vehicle repository
type VehicleRepository interface {
	// FindAll is a method that returns a map of all vehicles
	FindAll() (v map[int]Vehicle, err error)
	GetByRegistration(registration string) (v *Vehicle, err error)
	Create(vehicle VehicleAttributes) (v Vehicle, err error)
}
