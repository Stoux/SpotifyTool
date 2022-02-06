package state

import "os"

var (
	port    = "8080"
	apiRoot = ""
)

func ResolveAddress() {
	if resolvedPort := os.Getenv("SPOTIFY_TOOL_PORT"); resolvedPort != "" {
		port = resolvedPort
	}
	if apiRoot = os.Getenv("SPOTIFY_ROOL_PUBLIC_ROOT"); apiRoot == "" {
		apiRoot = "http://localhost:" + port
	}
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
