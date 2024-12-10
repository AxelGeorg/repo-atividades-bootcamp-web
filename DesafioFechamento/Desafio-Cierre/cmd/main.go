package main

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/repository/loader"
	"app/internal/service"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := &ConfigAppDefault{
		ServerAddr: os.Getenv("SERVER_ADDR"),
		DbFile:     os.Getenv("DB_FILE"),
	}

	app := NewApplicationDefault(cfg)

	err := app.SetUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type ConfigAppDefault struct {
	ServerAddr string
	DbFile     string
}

func NewApplicationDefault(cfg *ConfigAppDefault) *ApplicationDefault {
	defaultRouter := chi.NewRouter()
	defaultConfig := &ConfigAppDefault{
		ServerAddr: ":8080",
		DbFile:     "../docs/db/tickets.csv",
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

type ApplicationDefault struct {
	rt         *chi.Mux
	serverAddr string
	dbFile     string
}

func (a *ApplicationDefault) SetUp() (err error) {
	db := loader.NewLoaderTicketCSV(a.dbFile)
	tickets, err := db.Load()

	if err != nil {
		fmt.Printf("Error loading tickets: %v\n", err)
		return
	}

	rp := repository.NewRepositoryTicketMap(tickets)
	sv := service.NewServiceTicketDefault(&rp)
	hd := handler.NewHandlerTickets(&sv)

	(*a).rt.Get("/ticket/getByCountry/{dest}", hd.GetByCountry)
	(*a).rt.Get("/ticket/getAverage/{dest}", hd.GetAverage)

	return
}

func (a *ApplicationDefault) Run() (err error) {
	err = http.ListenAndServe(a.serverAddr, a.rt)
	return
}
