package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// both of these packages are fetched from github
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Remember a := App{}?
// This is that struct
type App struct {
	// Router is a mux Router, this will handle our handlers and prepare them for server
	Router *mux.Router
	// Broker is our custom struct that deals with database initalization and schema migration
	Broker Broker
}

// Initialize() initializes a connection to a database and handlers setup
func (a *App) Initialize() {
	// Same as a := App{}
	// But this way we can do more logic before initialization.
	// NewBroker() is a custom method defined in broker.go
	a.Broker = NewBroker()
	// Get all four env keys and assign them to variables.
	PgUsername := os.Getenv("PG_USERNAME")
	PgPassword := os.Getenv("PG_PASSWORD")
	PgDbName := os.Getenv("PG_DB_NAME")
	PgDbHost := os.Getenv("PG_DB_HOST")
	// You might be asking - "why do I need that?"
	// This is forward thinking - if I'd like to plug literally any other database to our application,
	// I could define an interface here, rename the function to SetConfig
	// and create the same function in our new broker.
	// This way the application is ready to extend functionality when needed.
	a.Broker.SetPostgresConfig(PgUsername, PgPassword, PgDbName, PgDbHost)
	// Same as in main.go, but this time we're checking for error.
	// In Go you'll see this a lot:
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	// This is a way how we deal with errors.
	if err := a.Broker.InitializeBroker(); err != nil {
		log.Fatalf("Error initializing postgres connection: %v", err)
	}

	// Create a router of type *mux.Router (imported type from Mux, referencing the value)
	router := mux.NewRouter()

	// Easy way to prepare our application for versioning.
	// Although Mux has it's own way of dealing with prefixes, I personally don't like it.
	// Hence my prefix and then concatenation.
	prefix := "/api"

	// Every handler will require the path for the endpoint, function that it'll call
	// and a method on which it'll be used
	router.Handle(prefix+"/post", a.GetAllPost()).Methods(http.MethodGet)
	// {post_id} is a parameter in the path. We'll grab it in the handler logic in handlers.go
	router.Handle(prefix+"/post/{post_id}", a.GetSinglePost()).Methods(http.MethodGet)
	// Why do you use `http.MethodPost` instead of typing "POST"?
	// Small preference - all methods are defined in net/http package so I'm extra careful.
	router.Handle(prefix+"/post", a.CreatePost()).Methods(http.MethodPost)
	// TODO: add handler for updating post:

	router.Handle(prefix+"/post/{post_id}", a.DeletePost()).Methods(http.MethodDelete)

	// Everything below is a helper function that will print all available routes in the api.
	// No need to explain it more.
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
	// Assign router to our application.
	a.Router = router
}

// Run() is the function that will run the server.
// Takes addr as a host - this way, if we need, we can define the port as env var.
func (a *App) Run(addr string) {
	// Very basic CORS setup. Better to have it now rather than wait for frustrated FrontEnd devs ;)
	handler := cors.Default().Handler(a.Router)
	log.Printf("Server is listening on %v", addr)
	// Run our server!
	http.ListenAndServe(addr, handler)
}
