package main

import (
	"SpotifyTool/processing"
	"SpotifyTool/server"
)

var shutdown chan bool

func main() {
	shutdown = make(chan bool)

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
