package server

import (
	"SpotifyTool/server/routes"
	"SpotifyTool/server/state"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	apiServer      *http.Server
	frontendServer *http.Server
)

func Serve(shutdown chan bool) {
	if apiServer != nil || frontendServer != nil {
		panic("Can only create an HTTP server once!")
	}

	api(shutdown)
	frontend(shutdown)
}

func api(shutdown chan bool) {
	// Build the route and it's routes
	router := mux.NewRouter()
	routes.GeneralRoutes(router)
	routes.AuthRoutes(router)

	// Build the server with those routes
	apiServer = &http.Server{
		Addr:         ":" + state.GetApiPort(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the goroutine that will handle all requests
	go func() {
		fmt.Println("Starting API server on", apiServer.Addr)
		if err := apiServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		shutdown <- true
	}()
}

func frontend(shutdown chan bool) {
	dirPath := state.GetFrontendServePath()
	if dirPath == "" {
		fmt.Println("No frontend path set / frontend serve disabled.")
		return
	}

	// Server will have a single route that's just the build assets folder
	dir := os.DirFS(dirPath)
	fileServer := http.FileServer(http.FS(dir))
	router := mux.NewRouter()
	router.Handle("/", fileServer)

	// Build the server
	frontendServer = &http.Server{
		Addr:         ":" + state.GetFrontendPort(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the goroutine that will handle all requests
	go func() {
		fmt.Println("Starting Frontend server on", frontendServer.Addr, "serving folder", dirPath)
		if err := frontendServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		shutdown <- true
	}()
}
