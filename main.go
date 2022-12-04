package main

import (
	"SpotifyTool/persistance"
	"SpotifyTool/processing"
	"SpotifyTool/server"
	httpState "SpotifyTool/server/state"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

var shutdown chan bool

func main() {
	shutdown = make(chan bool)

	initSentry()

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

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)
}
