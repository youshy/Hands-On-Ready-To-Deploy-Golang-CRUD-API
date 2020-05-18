package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	Broker Broker
}

func (a *App) Initialize() {
	a.Broker = NewBroker()
	PgUsername := os.Getenv("PG_USERNAME")
	PgPassword := os.Getenv("PG_PASSWORD")
	PgDbName := os.Getenv("PG_DB_NAME")
	PgDbHost := os.Getenv("PG_DB_HOST")
	a.Broker.SetPostgresConfig(PgUsername, PgPassword, PgDbName, PgDbHost)
	if err := a.Broker.InitializeBroker(); err != nil {
		log.Fatalf("Error initializing postgres connection: %v", err)
	}

	router := mux.NewRouter()

	prefix := "/api"

	router.Handle(prefix+"/post", a.GetAllPost()).Methods(http.MethodGet)
	router.Handle(prefix+"/post/{post_id}", a.GetSinglePost()).Methods(http.MethodGet)
	router.Handle(prefix+"/post", a.CreatePost()).Methods(http.MethodPost)
	// TODO: add handler for updating post:

	router.Handle(prefix+"/post/{post_id}", a.DeletePost()).Methods(http.MethodDelete)

	log.Printf("Available routes:\n")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", m, t)
		return nil
	})
	a.Router = router
}

func (a *App) Run(addr string) {
	handler := cors.Default().Handler(a.Router)
	log.Printf("Server is listening on %v", addr)
	http.ListenAndServe(addr, handler)
}
