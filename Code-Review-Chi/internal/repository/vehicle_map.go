package repository

import (
	"app/internal"
	"errors"
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

func (r *VehicleMap) GetById(id int) (*internal.Vehicle, error) {
	vehicle, ok := r.db[id]
	if !ok {
		return nil, errors.New("erro")
	}

	return &vehicle, nil
}

func (r *VehicleMap) GetByRegistration(registration string) (*internal.Vehicle, error) {
	for _, value := range r.db {
		if value.Registration == registration {
			return &value, nil
		}
	}

	return nil, nil
}

func (r *VehicleMap) Patch(id int, updates map[string]interface{}) (*internal.Vehicle, error) {
	vehicle, err := r.GetById(id)
	if err != nil {
		return nil, err
	}

	if maxSpeed, ok := updates["max_speed"].(float64); ok {
		vehicle.MaxSpeed = maxSpeed
	}

	r.db[id] = *vehicle

	return vehicle, nil
}
