package main

import (
	"SpotifyTool/server"
	"time"
)

func main() {

	server.Serve()

	// Do other things
	for true {
		time.Sleep(1)
	}
}
