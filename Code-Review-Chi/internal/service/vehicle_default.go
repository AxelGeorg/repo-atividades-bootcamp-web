package service

import (
	"app/internal"
	"app/internal/errors"
)

// NewVehicleDefault is a function that returns a new instance of VehicleDefault
func NewVehicleDefault(rp internal.VehicleRepository) VehicleDefault {
	return VehicleDefault{rp: rp}
}

// VehicleDefault is a struct that represents the default service for vehicles
type VehicleDefault struct {
	// rp is the repository that will be used by the service
	rp internal.VehicleRepository
}

// FindAll is a method that returns a map of all vehicles
func (s *VehicleDefault) FindAll() (v map[int]internal.Vehicle, err error) {
	v, err = s.rp.FindAll()
	return
}

func validateVehicle(repo internal.VehicleRepository, vehicle internal.VehicleAttributes) error {
	if vehicle.Brand == "" {
		return errors.NewBadRequestError("brand is required")
	}
	if vehicle.Model == "" {
		return errors.NewBadRequestError("model is required")
	}
	if vehicle.Registration == "" {
		return errors.NewBadRequestError("registration is required")
	}
	if vehicle.Color == "" {
		return errors.NewBadRequestError("color is required")
	}
	if vehicle.FabricationYear <= 0 {
		return errors.NewBadRequestError("fabrication year must be a positive value")
	}
	if vehicle.Capacity <= 0 {
		return errors.NewBadRequestError("capacity must be greater than zero")
	}
	if vehicle.MaxSpeed < 0 {
		return errors.NewBadRequestError("max speed cannot be negative")
	}
	if vehicle.FuelType == "" {
		return errors.NewBadRequestError("fuel type is required")
	}
	if vehicle.Transmission == "" {
		return errors.NewBadRequestError("transmission is required")
	}
	if vehicle.Weight <= 0 {
		return errors.NewBadRequestError("weight must be a positive value")
	}
	if vehicle.Dimensions.Length <= 0 || vehicle.Dimensions.Width <= 0 || vehicle.Dimensions.Height <= 0 {
		return errors.NewBadRequestError("dimensions must be positive values")
	}

	v, err := repo.GetByRegistration(vehicle.Registration)
	if err != nil {
		return err
	}

	if v != nil {
		return errors.NewConflictError("vehicle with this registration already exists")
	}

	return nil
}

func (s *VehicleDefault) Create(vehicle internal.VehicleAttributes) (v internal.Vehicle, err error) {

	err = validateVehicle(s.rp, vehicle)
	if err != nil {
		return
	}

	v, err = s.rp.Create(vehicle)
	return
}
