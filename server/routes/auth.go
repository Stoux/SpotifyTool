package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/processing"
	"SpotifyTool/server/handlers"
	httpState "SpotifyTool/server/state"
	"SpotifyTool/server/util"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	netUrl "net/url"
)

var (
	knownOAuthStates                            = map[string]string{}
	authenticator    *spotifyauth.Authenticator = nil
)

func AuthRoutes(router *mux.Router) {
	authenticator = httpState.GetSpotifyAuthenticator()

	// Register the route
	router.HandleFunc("/auth/start", handlers.JsonWithOutput(handleAuthStart))
	router.HandleFunc("/auth/callback", handleCallback)
	router.HandleFunc("/auth/me", handlers.Auth(getCurrentUser))
}

func handleAuthStart(writer http.ResponseWriter, request *http.Request) interface{} {
	// Create and store a state
	state := util.GetRandomString(16)
	knownOAuthStates[state] = state

	// Generate the auth URL
	authUrl := authenticator.AuthURL(state)

	// TODO: Queue the state for removal after 15 minutes

	return authStartResponse{authUrl}
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Fetch the state & check if it's known
	state := r.URL.Query().Get("state")
	if state == "" || knownOAuthStates[state] == "" {
		incorrectCallback(w, r, "Invalid OAuth State")
		return
	}

	// State has been 'used', remove it
	delete(knownOAuthStates, state)

	// Fetch the token from Spotify
	spotifyToken, err := authenticator.Token(r.Context(), state, r)
	if err != nil {
		incorrectCallback(w, r, "Unable to fetch token")
		return
	}

	// Build a client & fetch the user
	client := spotify.New(authenticator.Client(r.Context(), spotifyToken))
	user, err := client.CurrentUser(r.Context())
	if err != nil {
		incorrectCallback(w, r, "Successfully authenticated but failed to fetch user data?")
		return
	}

	// Pass the user to processing for account creation & further handling
	log.Println("User has authenticated: " + user.DisplayName + "(" + user.ID + ")")
	processing.GetLoginChannel() <- processing.SpotifyClientLogin{
		Token: spotifyToken,
		User:  user,
	}

	// Create an access-token
	accessToken := httpState.CreateTokenFor(user.ID)

	redirectUrl := httpState.GetFrontendRoot() + "/authenticated?" +
		"token=" + netUrl.QueryEscape(accessToken.Identifier) +
		"&name=" + netUrl.QueryEscape(user.DisplayName)
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func getCurrentUser(w http.ResponseWriter, r *http.Request, userSpotifyId string) {
	handlers.JsonWithOutput(func(w http.ResponseWriter, r *http.Request) (result interface{}) {
		user := models.ToolUser{}
		persistance.Db.Where("spotify_id = ?", userSpotifyId).Find(&user)

		return user
	})(w, r)
}

func incorrectCallback(writer http.ResponseWriter, request *http.Request, error string) {
	redirectUrl := httpState.GetFrontendRoot() + "/login?error=" + netUrl.QueryEscape(error)
	http.Redirect(writer, request, redirectUrl, http.StatusFound)
}

type authStartResponse struct {
	Url string `json:"url"`
}
