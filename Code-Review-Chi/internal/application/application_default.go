package application

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress string
	// LoaderFilePath is the path to the file that contains the vehicles
	LoaderFilePath string
}

// NewServerChi is a function that returns a new instance of ServerChi
func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	// default values
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
		if cfg.LoaderFilePath != "" {
			defaultConfig.LoaderFilePath = cfg.LoaderFilePath
		}
	}

	return &ServerChi{
		serverAddress:  defaultConfig.ServerAddress,
		loaderFilePath: defaultConfig.LoaderFilePath,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress string
	// loaderFilePath is the path to the file that contains the vehicles
	loaderFilePath string
}

// Run is a method that runs the application
func (a *ServerChi) Run() (err error) {
	// dependencies
	// - loader
	ld := loader.NewVehicleJSONFile(a.loaderFilePath)
	db, err := ld.Load()
	if err != nil {
		return
	}
	// - repository
	rp := repository.NewVehicleMap(db)
	// - service
	sv := service.NewVehicleDefault(&rp)
	// - handler
	hd := handler.NewVehicleDefault(&sv)
	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints
	rt.Route("/vehicles", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Get("/color/{color}/year/{year}", hd.GetColorYear)
		r.Get("/brand/{brand}/between/{start_year}/{end_year}", hd.GetBrandAndYearsPeriod)
		r.Get("/average_speed/brand/{brand}", hd.GetAverageSpeed)
		r.Post("/batch/", hd.PostMany)
		r.Post("/", hd.Post)
		r.Put("/{id}/update_speed", hd.PutSpeed)
		r.Get("/fuel_type/{type}", hd.GetFuelType)
		r.Delete("/{id}", hd.Delete)
		r.Get("/transmission/{type}", hd.GetTransmission)

		r.Put("/{id}/update_fuel", hd.PutFuel)
		r.Get("/average_capacity/brand/{brand}", hd.GetAverageBrand)
		r.Get("/dimensions?length={min_length}-{max_length}&width={min_width}-{max_width}", hd.GetDimensions)
		r.Get("/weight?min={weight_min}&max={weight_max}", hd.GetWeight)
	})

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
