package server

import (
	"SpotifyTool/server/routes"
	"SpotifyTool/server/state"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func Serve(shutdown chan bool) {
	if server != nil {
		panic("Can only create an HTTP server once!")
	}

	state.ResolveAddress()

	// Build the route and it's routes
	router := mux.NewRouter()
	routes.GeneralRoutes(router)
	routes.AuthRoutes(router)

	// Build the server with those routes
	server = &http.Server{
		Addr:         ":" + state.GetBindPort(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the goroutine that will handle all requests
	go func() {
		fmt.Println("Starting server on ", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		shutdown <- true
	}()
}
