package main

import (
	"SpotifyTool/persistance"
	"SpotifyTool/processing"
	"SpotifyTool/server"
)

var shutdown chan bool

func main() {
	shutdown = make(chan bool)

	// Create the database (connection)
	persistance.Init()

	// Start the background processing
	processing.Init()

	// Start the http server
	server.Serve(shutdown)

	// Do other things
	for {
		if shouldShutdown := <-shutdown; shouldShutdown {
			return
		}
	}
}
