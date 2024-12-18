package handler

import (
	"app/internal"
	errorss "app/internal/errors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bootcamp-go/web/response"
)

const (
	MessageVehicleCreated = "Vehicle created"
	MessageVehicleUpdated = "Vehicle updated"
	MessageVehicleDeleted = "Vehicle deleted"
)

// VehicleJSON is a struct that represents a vehicle in JSON format
type VehicleJSON struct {
	ID              int     `json:"id"`
	Brand           string  `json:"brand"`
	Model           string  `json:"model"`
	Registration    string  `json:"registration"`
	Color           string  `json:"color"`
	FabricationYear int     `json:"year"`
	Capacity        int     `json:"passengers"`
	MaxSpeed        float64 `json:"max_speed"`
	FuelType        string  `json:"fuel_type"`
	Transmission    string  `json:"transmission"`
	Weight          float64 `json:"weight"`
	Height          float64 `json:"height"`
	Length          float64 `json:"length"`
	Width           float64 `json:"width"`
}

type RequestBodyFuelType struct {
	FuelType string `json:"fuel_type"`
}

type RequestBodyVehicle struct {
	Brand           string  `json:"brand"`
	Model           string  `json:"model"`
	Registration    string  `json:"registration"`
	Color           string  `json:"color"`
	FabricationYear int     `json:"year"`
	Capacity        int     `json:"passengers"`
	MaxSpeed        float64 `json:"max_speed"`
	FuelType        string  `json:"fuel_type"`
	Transmission    string  `json:"transmission"`
	Weight          float64 `json:"weight"`
	Height          float64 `json:"height"`
	Length          float64 `json:"length"`
	Width           float64 `json:"width"`
}

type ResponseBodyVehicle struct {
	Message string              `json:"message"`
	Data    *RequestBodyVehicle `json:"data,omitempty"`
	Error   bool                `json:"error"`
}

func ResponseWithError(w http.ResponseWriter, err error, statusCode int) {
	body := &ResponseBodyVehicle{
		Message: http.StatusText(statusCode) + " - " + err.Error(),
		Data:    nil,
		Error:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func RespondWithVehicle(w http.ResponseWriter, vehicle *internal.Vehicle, statusCode int, message string) {
	var body *ResponseBodyVehicle
	if vehicle == nil {
		body = &ResponseBodyVehicle{
			Message: message,
			Data:    nil,
			Error:   false,
		}
	} else {
		dt := RequestBodyVehicle{
			Brand:           vehicle.Brand,
			Model:           vehicle.Model,
			Registration:    vehicle.Registration,
			Color:           vehicle.Color,
			FabricationYear: vehicle.FabricationYear,
			Capacity:        vehicle.Capacity,
			MaxSpeed:        vehicle.MaxSpeed,
			FuelType:        vehicle.FuelType,
			Transmission:    vehicle.Transmission,
			Weight:          vehicle.Weight,
			Height:          vehicle.Height,
			Length:          vehicle.Length,
			Width:           vehicle.Width,
		}

		body = &ResponseBodyVehicle{
			Message: message,
			Data:    &dt,
			Error:   false,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

// NewVehicleDefault is a function that returns a new instance of VehicleDefault
func NewVehicleDefault(sv internal.VehicleService) *VehicleDefault {
	return &VehicleDefault{sv: sv}
}

// VehicleDefault is a struct with methods that represent handlers for vehicles
type VehicleDefault struct {
	// sv is the service that will be used by the handler
	sv internal.VehicleService
}

// GetAll is a method that returns a handler for the route GET /vehicles
func (h *VehicleDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	// request
	// ...

	// process
	// - get all vehicles
	v, err := h.sv.FindAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	// response
	data := make(map[int]VehicleJSON)
	for key, value := range v {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "success",
		"data":    data,
	})
}

func (h *VehicleDefault) GetColorYear(w http.ResponseWriter, r *http.Request) {
	colorAndYear := r.URL.Path[len("/vehicles/color/"):]

	params := strings.Split(colorAndYear, "/")
	if len(params) != 3 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	fabricationYear, err := strconv.Atoi(params[2])
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	filter := internal.VehicleAttributesFilter{
		Color:                params[0],
		FabricationYearStart: fabricationYear,
		FabricationYearEnd:   fabricationYear,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *VehicleDefault) GetBrandAndYearsPeriod(w http.ResponseWriter, r *http.Request) {
	brancAndYearsPeriod := r.URL.Path[len("/vehicles/brand/"):]

	params := strings.Split(brancAndYearsPeriod, "/")
	if len(params) != 4 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	yearStart, err := strconv.Atoi(params[2])
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	yearEnd, err := strconv.Atoi(params[3])
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	filter := internal.VehicleAttributesFilter{
		Brand:                params[0],
		FabricationYearStart: yearStart,
		FabricationYearEnd:   yearEnd,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *VehicleDefault) GetAverageSpeed(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Path[len("/vehicles/average_speed/brand/"):]
	avarage, err := h.sv.GetAverageSpeed(brand)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprint("Avarage Speed: ", avarage))
}

func (h *VehicleDefault) Post(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBodyVehicle
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	dimensions := internal.Dimensions{
		Height: reqBody.Height,
		Length: reqBody.Length,
		Width:  reqBody.Width,
	}

	vehicle := internal.VehicleAttributes{
		Brand:           reqBody.Brand,
		Model:           reqBody.Model,
		Registration:    reqBody.Registration,
		Color:           reqBody.Color,
		FabricationYear: reqBody.FabricationYear,
		Capacity:        reqBody.Capacity,
		MaxSpeed:        reqBody.MaxSpeed,
		FuelType:        reqBody.FuelType,
		Transmission:    reqBody.Transmission,
		Weight:          reqBody.Weight,
		Dimensions:      dimensions,
	}

	productServ, err := h.sv.Create(vehicle)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	RespondWithVehicle(w, &productServ, http.StatusCreated, MessageVehicleCreated)
}

func (h *VehicleDefault) PostMany(w http.ResponseWriter, r *http.Request) {
	var reqBodies []RequestBodyVehicle
	if err := json.NewDecoder(r.Body).Decode(&reqBodies); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	for _, reqBody := range reqBodies {
		dimensions := internal.Dimensions{
			Height: reqBody.Height,
			Length: reqBody.Length,
			Width:  reqBody.Width,
		}

		vehicle := internal.VehicleAttributes{
			Brand:           reqBody.Brand,
			Model:           reqBody.Model,
			Registration:    reqBody.Registration,
			Color:           reqBody.Color,
			FabricationYear: reqBody.FabricationYear,
			Capacity:        reqBody.Capacity,
			MaxSpeed:        reqBody.MaxSpeed,
			FuelType:        reqBody.FuelType,
			Transmission:    reqBody.Transmission,
			Weight:          reqBody.Weight,
			Dimensions:      dimensions,
		}

		_, err := h.sv.Create(vehicle)
		if err != nil {
			ResponseWithError(w, err, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Vehicles created successfully")
}

func (h *VehicleDefault) PutSpeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/vehicles/"):]

	params := strings.Split(url, "/")
	if len(params) != 2 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	idVehicle, err := strconv.Atoi(params[0])
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	vehicle, err := h.sv.Patch(idVehicle, updates)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

func (h *VehicleDefault) GetFuelType(w http.ResponseWriter, r *http.Request) {
	fuelType := r.URL.Path[len("/vehicles/fuel_type/"):]
	if fuelType == "" {
		ResponseWithError(w, errors.New("erro"), http.StatusBadRequest)
		return
	}

	filter := internal.VehicleAttributesFilter{
		FuelType: fuelType,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *VehicleDefault) Delete(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/vehicles/"):]
	if url == "" {
		ResponseWithError(w, errors.New("erro"), http.StatusBadRequest)
		return
	}

	idVehicle, err := strconv.Atoi(url)
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	err = h.sv.Delete(idVehicle)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *VehicleDefault) GetTransmission(w http.ResponseWriter, r *http.Request) {
	transmission := r.URL.Path[len("/vehicles/transmission/"):]
	if transmission == "" {
		ResponseWithError(w, errors.New("erro"), http.StatusBadRequest)
		return
	}

	filter := internal.VehicleAttributesFilter{
		Transmission: transmission,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *VehicleDefault) PutFuel(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/vehicles/"):]
	params := strings.Split(url, "/")
	if len(params) != 2 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	idVehicle, err := strconv.Atoi(params[0])
	if err != nil {
		ResponseWithError(w, errors.New("fabrication year must be a valid integer"), http.StatusBadRequest)
		return
	}

	var update RequestBodyFuelType
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	vehicle, err := h.sv.PutFuel(idVehicle, update.FuelType)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

func (h *VehicleDefault) GetAverageCapacity(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Path[len("/vehicles/average_capacity/brand/"):]
	avarage, err := h.sv.GetAverageCapacity(brand)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprint("Avarage Capacity: ", avarage))
}

func (h *VehicleDefault) GetDimensions(w http.ResponseWriter, r *http.Request) {
	//search
	//url := r.URL.Path[len("/vehicles/dimensions?length={min_length}-{max_length}&width={min_width}-{max_width}"):]

	paramLength := r.URL.Query().Get("length")
	paramWidth := r.URL.Query().Get("width")

	paramsLength := strings.Split(paramLength, "-")
	if len(paramsLength) != 2 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	paramsWidth := strings.Split(paramWidth, "-")
	if len(paramsLength) != 2 {
		ResponseWithError(w, errors.New("the URL is not in the correct format"), http.StatusBadRequest)
		return
	}

	lendthMin, err := strconv.ParseFloat(paramsLength[0], 64)
	if err != nil {
		ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	lendthMax, err := strconv.ParseFloat(paramsLength[1], 64)
	if err != nil {
		ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	widthMin, err := strconv.ParseFloat(paramsWidth[0], 64)
	if err != nil {
		ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	widthMax, err := strconv.ParseFloat(paramsWidth[1], 64)
	if err != nil {
		ResponseWithError(w, errors.New("invalid price format"), http.StatusBadRequest)
		return
	}

	dimMin := internal.Dimensions{
		Length: lendthMin,
		Width:  widthMin,
	}

	dimMax := internal.Dimensions{
		Length: lendthMax,
		Width:  widthMax,
	}

	filter := internal.VehicleAttributesFilter{
		DimensionMin: dimMin,
		DimensionMax: dimMax,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *VehicleDefault) GetWeight(w http.ResponseWriter, r *http.Request) {
	//search
	//url := r.URL.Path[len("/vehicles/weight?min={weight_min}&max={weight_max}"):]

	weightMinStr := r.URL.Query().Get("min")
	weightMaxStr := r.URL.Query().Get("max")

	weightMin, err := strconv.ParseFloat(weightMinStr, 64)
	if err != nil {
		ResponseWithError(w, errors.New("formato de peso mínimo inválido"), http.StatusBadRequest)
		return
	}

	weightMax, err := strconv.ParseFloat(weightMaxStr, 64)
	if err != nil {
		ResponseWithError(w, errors.New("formato de peso máximo inválido"), http.StatusBadRequest)
		return
	}

	filter := internal.VehicleAttributesFilter{
		WeightMin: weightMin,
		WeightMax: weightMax,
	}

	vehicles, err := h.sv.GetVehiclesWithFilter(filter)
	if customErr, ok := err.(*errorss.CustomError); ok {
		http.Error(w, customErr.Message, customErr.StatusHttp)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := make(map[int]VehicleJSON)
	for key, value := range *vehicles {
		data[key] = VehicleJSON{
			ID:              value.Id,
			Brand:           value.Brand,
			Model:           value.Model,
			Registration:    value.Registration,
			Color:           value.Color,
			FabricationYear: value.FabricationYear,
			Capacity:        value.Capacity,
			MaxSpeed:        value.MaxSpeed,
			FuelType:        value.FuelType,
			Transmission:    value.Transmission,
			Weight:          value.Weight,
			Height:          value.Height,
			Length:          value.Length,
			Width:           value.Width,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
