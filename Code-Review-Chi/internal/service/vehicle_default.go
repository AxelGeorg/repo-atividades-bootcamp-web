package service

import (
	"app/internal"
	errorss "app/internal/errors"
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
		return errorss.NewBadRequestError("brand is required")
	}
	if vehicle.Model == "" {
		return errorss.NewBadRequestError("model is required")
	}
	if vehicle.Registration == "" {
		return errorss.NewBadRequestError("registration is required")
	}
	if vehicle.Color == "" {
		return errorss.NewBadRequestError("color is required")
	}
	if vehicle.FabricationYear <= 0 {
		return errorss.NewBadRequestError("fabrication year must be a positive value")
	}
	if vehicle.Capacity <= 0 {
		return errorss.NewBadRequestError("capacity must be greater than zero")
	}
	if vehicle.MaxSpeed < 0 {
		return errorss.NewBadRequestError("max speed cannot be negative")
	}
	if vehicle.FuelType == "" {
		return errorss.NewBadRequestError("fuel type is required")
	}
	if vehicle.Transmission == "" {
		return errorss.NewBadRequestError("transmission is required")
	}
	if vehicle.Weight <= 0 {
		return errorss.NewBadRequestError("weight must be a positive value")
	}
	if vehicle.Dimensions.Length <= 0 || vehicle.Dimensions.Width <= 0 || vehicle.Dimensions.Height <= 0 {
		return errorss.NewBadRequestError("dimensions must be positive values")
	}

	v, err := repo.GetByRegistration(vehicle.Registration)
	if err != nil {
		return err
	}

	if v != nil {
		return errorss.NewConflictError("vehicle with this registration already exists")
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

func (s *VehicleDefault) GetAverageSpeed(brand string) (float64, error) {

	filter := internal.VehicleAttributesFilter{
		Brand: brand,
	}

	vehicles, err := s.GetVehiclesWithFilter(filter)
	if err != nil {
		return 0.0, nil
	}

	var sumMaxSpeed float64
	for _, vehicle := range *vehicles {
		sumMaxSpeed += vehicle.MaxSpeed
	}

	return (sumMaxSpeed / float64(len(*vehicles))), nil
}

func (s *VehicleDefault) GetVehiclesWithFilter(filter internal.VehicleAttributesFilter) (*map[int]internal.Vehicle, error) {
	list, err := s.rp.FindAll()
	if err != nil {
		return nil, errorss.NewBadRequestError("no vehicles")
	}

	mapReturn := make(map[int]internal.Vehicle)

	for _, vehicle := range list {
		if filter.Brand != "" {
			if filter.Brand != vehicle.Brand {
				continue
			}
		}

		if filter.Model != "" {
			if filter.Model != vehicle.Model {
				continue
			}
		}

		if filter.Registration != "" {
			if filter.Registration != vehicle.Registration {
				continue
			}
		}

		if filter.Color != "" {
			if filter.Color != vehicle.Color {
				continue
			}
		}

		if filter.FabricationYearStart > 0 && filter.FabricationYearStart == filter.FabricationYearEnd {
			if filter.FabricationYearStart != vehicle.FabricationYear {
				continue
			}
		} else if filter.FabricationYearStart > 0 && filter.FabricationYearEnd > 0 {
			if vehicle.FabricationYear > filter.FabricationYearEnd || vehicle.FabricationYear < filter.FabricationYearStart {
				continue
			}
		}

		if filter.Capacity > 0 {
			if filter.Capacity != vehicle.Capacity {
				continue
			}
		}

		if filter.MaxSpeed > 0 {
			if filter.MaxSpeed != vehicle.MaxSpeed {
				continue
			}
		}

		if filter.FuelType != "" {
			if filter.FuelType != vehicle.FuelType {
				continue
			}
		}

		if filter.Transmission != "" {
			if filter.Transmission != vehicle.Transmission {
				continue
			}
		}

		if filter.Weight > 0 {
			if filter.Weight != vehicle.Weight {
				continue
			}
		}

		if filter.Dimensions.Length > 0 {
			if filter.Dimensions.Length != vehicle.Dimensions.Length {
				continue
			}
		}

		if filter.Dimensions.Width > 0 {
			if filter.Dimensions.Width != vehicle.Dimensions.Width {
				continue
			}
		}

		if filter.Dimensions.Height > 0 {
			if filter.Dimensions.Height != vehicle.Dimensions.Height {
				continue
			}
		}

		mapReturn[vehicle.Id] = vehicle
	}

	return &mapReturn, nil
}

func (s *VehicleDefault) Patch(id int, updates map[string]interface{}) (*internal.Vehicle, error) {
	vehicle, err := s.rp.Patch(id, updates)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *VehicleDefault) Delete(id int) error {
	err := s.rp.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *VehicleDefault) PutFuel(id int, fuelType string) (*internal.Vehicle, error) {
	vehicle, err := s.rp.PutFuel(id, fuelType)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}
