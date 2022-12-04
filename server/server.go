package server

import (
	"SpotifyTool/server/handlers"
	"SpotifyTool/server/routes"
	"SpotifyTool/server/state"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	routes.PlaylistRoutes(router)
	routes.PlaylistBackupRoutes(router)
	routes.TrackRoutes(router)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			next.ServeHTTP(w, r)
		})
	})

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
	router := mux.NewRouter()
	router.HandleFunc("/api/health", handlers.JsonWithOutput(func(w http.ResponseWriter, r *http.Request) (result interface{}) {
		return map[string]bool{"ok": true}
	}))
	router.PathPrefix("/").Handler(
		handlers.SpaHandler{
			StaticPath: dirPath,
			IndexPath:  "index.html",
		},
	)

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
