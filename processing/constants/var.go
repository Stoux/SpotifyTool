package constants

import "time"

const (
	MaxPlaylistsPerPage                = 50
	MaxTracksPerPage                   = 100
	FetchPlaylistsInterval             = 15 * time.Minute
	FetchSpotifyOwnedPlaylistsInterval = 2 * time.Hour
)
