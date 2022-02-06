package state

import (
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"os"
)

var (
	port                                     = "8080"
	apiRoot                                  = ""
	authenticator *spotifyauth.Authenticator = nil
)

func ResolveAddress() {
	if resolvedPort := os.Getenv("SPOTIFY_TOOL_PORT"); resolvedPort != "" {
		port = resolvedPort
	}
	if apiRoot = os.Getenv("SPOTIFY_ROOL_PUBLIC_ROOT"); apiRoot == "" {
		apiRoot = "http://localhost:" + port
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
	// TODO
	return "https://stoux.nl"
}

// GetApiRoot returns the root domain / HTTP address without trailing slash
func GetApiRoot() string {
	return apiRoot
}

func GetBindPort() string {
	return port
}

func GetSpotifyAuthenticator() *spotifyauth.Authenticator {
	return authenticator
}
