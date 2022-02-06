package main

import (
	"SpotifyTool/persistance"
	"SpotifyTool/processing"
	"SpotifyTool/server"
	httpState "SpotifyTool/server/state"
)

var shutdown chan bool

func main() {
	shutdown = make(chan bool)

	// Resolve any settings & address values
	httpState.ResolveAddress()

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
