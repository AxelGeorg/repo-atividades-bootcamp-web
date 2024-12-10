package repository

import (
	"app/internal"
	"fmt"
)

// NewVehicleMap is a function that returns a new instance of VehicleMap
func NewVehicleMap(dbMap map[int]internal.Vehicle) VehicleMap {
	// default db
	defaultDb := make(map[int]internal.Vehicle)
	if dbMap != nil {
		defaultDb = dbMap
	}
	return VehicleMap{db: defaultDb}
}

// VehicleMap is a struct that represents a vehicle repository
type VehicleMap struct {
	// db is a map of vehicles
	db map[int]internal.Vehicle
}

// FindAll is a method that returns a map of all vehicles
func (r *VehicleMap) FindAll() (v map[int]internal.Vehicle, err error) {
	v = make(map[int]internal.Vehicle)

	// copy db
	for key, value := range r.db {
		v[key] = value
	}

	return
}

func (r *VehicleMap) Create(vehicle internal.VehicleAttributes) (v internal.Vehicle, err error) {
	newId := len(r.db) + 1
	v = internal.Vehicle{Id: newId, VehicleAttributes: vehicle}

	r.db[newId] = v
	fmt.Println(r.db[newId])
	fmt.Println(r.db)
	return
}

func (r *VehicleMap) GetByRegistration(registration string) (*internal.Vehicle, error) {
	for _, value := range r.db {
		if value.Registration == registration {
			return &value, nil
		}
	}

	return nil, nil
}
