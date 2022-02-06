package routes

import (
	"SpotifyTool/server/handlers"
	"github.com/gorilla/mux"
	netHttp "net/http"
)

func GeneralRoutes(router *mux.Router) {

	router.HandleFunc("/", handlers.Json(func(writer netHttp.ResponseWriter, request *netHttp.Request) (result interface{}, shouldOutput bool) {
		return homeResponse{
			Name:        "Stoux' Spotify Tool API",
			Description: "A Spotify tool that offers a bunch of utility functions regarding playlists & more",
			Version:     "1.0.0", // TODO: Get from Build
		}, true
	}))

}

type homeResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}
