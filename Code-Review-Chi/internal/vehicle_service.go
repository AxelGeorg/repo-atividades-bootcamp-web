package internal

// VehicleService is an interface that represents a vehicle service
type VehicleService interface {
	// FindAll is a method that returns a map of all vehicles
	FindAll() (v map[int]Vehicle, err error)
	Create(vehicle VehicleAttributes) (v Vehicle, err error)
	GetVehiclesWithFilter(filter VehicleAttributesFilter) (*map[int]Vehicle, error)
	GetAverageSpeed(brand string) (float64, error)
	Patch(id int, updates map[string]interface{}) (*Vehicle, error)
}
