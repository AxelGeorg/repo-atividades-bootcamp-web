package main

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/repository/loader"
	"app/internal/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// env
	// ...

	// application
	// - config
	/*cfg := &ConfigAppDefault{
		ServerAddr: os.Getenv("SERVER_ADDR"),
		DbFile:     os.Getenv("DB_FILE"),
	}*/

	cfg := &ConfigAppDefault{
		//ServerAddr: os.Getenv("SERVER_ADDR"),
		DbFile: "../../docs/db/tickets.csv",
	}

	app := NewApplicationDefault(cfg)

	// - setup
	err := app.SetUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	// - run
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// ConfigAppDefault represents the configuration of the default application
type ConfigAppDefault struct {
	// serverAddr represents the address of the server
	ServerAddr string
	// dbFile represents the path to the database file
	DbFile string
}

// NewApplicationDefault creates a new default application
func NewApplicationDefault(cfg *ConfigAppDefault) *ApplicationDefault {
	// default values
	defaultRouter := chi.NewRouter()
	defaultConfig := &ConfigAppDefault{
		ServerAddr: ":8080",
		DbFile:     "",
	}
	if cfg != nil {
		if cfg.ServerAddr != "" {
			defaultConfig.ServerAddr = cfg.ServerAddr
		}
		if cfg.DbFile != "" {
			defaultConfig.DbFile = cfg.DbFile
		}
	}

	return &ApplicationDefault{
		rt:         defaultRouter,
		serverAddr: defaultConfig.ServerAddr,
		dbFile:     defaultConfig.DbFile,
	}
}

// ApplicationDefault represents the default application
type ApplicationDefault struct {
	// router represents the router of the application
	rt *chi.Mux
	// serverAddr represents the address of the server
	serverAddr string
	// dbFile represents the path to the database file
	dbFile string
}

// SetUp sets up the application
func (a *ApplicationDefault) SetUp() (err error) {

	// dependencies
	db := loader.NewLoaderTicketCSV(a.dbFile)
	tickets, err := db.Load()
	if err != nil {
		fmt.Printf("Error loading tickets: %v\n", err)
		return
	}

	rp := repository.NewRepositoryTicketMap(tickets)
	sv := service.NewServiceTicketDefault(&rp)
	hd := handler.NewHandlerTickets(&sv)

	// routes
	(*a).rt.Get("/ticket/getByCountry/", hd.GetByCountry)
	(*a).rt.Get("/ticket/getAverage/", hd.GetAverage)

	return
}

// Run runs the application
func (a *ApplicationDefault) Run() (err error) {
	err = http.ListenAndServe(a.serverAddr, a.rt)
	return
}
