package state

import (
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"os"
)

var (
	apiPort                                  = "8080"
	frontendPort                             = "8040"
	apiRoot                                  = ""
	frontendRoot                             = ""
	frontendPath                             = "./web/build"
	serveFrontend                            = true
	authenticator *spotifyauth.Authenticator = nil
)

func ResolveAddress() {
	if resolvedPort := os.Getenv("SPOTIFY_TOOL_API_PORT"); resolvedPort != "" {
		apiPort = resolvedPort
	}
	if resolvedPort := os.Getenv("SPOTIFY_TOOL_FRONTEND_PORT"); resolvedPort != "" {
		frontendPort = resolvedPort
	}
	if apiRoot = os.Getenv("SPOTIFY_TOOL_API_ROOT"); apiRoot == "" {
		apiRoot = "http://localhost:" + apiPort
	}
	if frontendRoot = os.Getenv("SPOTIFY_TOOL_FRONTEND_ROOT"); frontendRoot == "" {
		frontendRoot = "http://localhost:" + frontendPort
	}
	if newFrontendPath := os.Getenv("SPOTIFY_TOOL_FRONTEND_PATH"); newFrontendPath != "" {
		frontendPath = newFrontendPath
	}
	if shouldServeFrontend := os.Getenv("SPOTIFY_TOOL_SERVE_FRONTEND"); shouldServeFrontend != "" {
		serveFrontend = shouldServeFrontend == "true" || shouldServeFrontend == "1" || shouldServeFrontend == "yes"
	}

	spotifyRedirectUri := apiRoot + "/auth/callback"
	authenticator = spotifyauth.New(spotifyauth.WithRedirectURL(spotifyRedirectUri), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserLibraryModify,
		spotifyauth.ScopeUserTopRead,
	))
}

// GetFrontendRoot returns the root domain / HTTP address without trailing slash
func GetFrontendRoot() string {
	return frontendRoot
}

// GetFrontendPort on which the frontend / SPA server should bind
func GetFrontendPort() string {
	return frontendPort
}

// GetFrontendServePath from which path the frontend HTML files should be served
func GetFrontendServePath() string {
	if serveFrontend {
		return frontendPath
	} else {
		return ""
	}
}

// GetApiRoot returns the root domain / HTTP address without trailing slash
func GetApiRoot() string {
	return apiRoot
}

func GetApiPort() string {
	return apiPort
}

func GetSpotifyAuthenticator() *spotifyauth.Authenticator {
	return authenticator
}
