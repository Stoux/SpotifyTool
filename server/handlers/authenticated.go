package handlers

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/state"
	"net/http"
)

// Auth validates the user is using a valid & known token and passes the resolved user's spotify ID to the wrapped handler.
func Auth(handler func(writer http.ResponseWriter, request *http.Request, userSpotifyId string)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "OPTIONS" {
			return
		}

		// Check if a valid Auth header was passed
		if identifier := request.Header.Get("Authorization"); identifier != "" {
			token := state.GetTokenBy(identifier)
			if token != nil {
				handler(writer, request, token.SpotifyId)
				return
			}
		}

		// Otherwise, we're going to return a 403
		writer.WriteHeader(http.StatusForbidden)
	}
}

func AuthUser(handler func(writer http.ResponseWriter, request *http.Request, user models.ToolUser)) func(writer http.ResponseWriter, request *http.Request) {
	return Auth(func(writer http.ResponseWriter, request *http.Request, userSpotifyId string) {
		handler(writer, request, persistance.GetToolUserBySpotifyId(userSpotifyId))
	})
}
