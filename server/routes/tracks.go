package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/handlers"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"time"
)

func TrackRoutes(router *mux.Router) {
	router.HandleFunc("/tracks/search", handlers.AuthUser(searchTracks))
}

type SearchFoundTrack struct {
	PlaylistId   string
	PlaylistName string
	TrackId      string
	TrackName    string
	Artists      string
	Album        string
	AddedAt      time.Time
	DeletedAt    gorm.DeletedAt
}

func searchTracks(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	// Get the query argument
	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, "Missing query argument", http.StatusBadRequest)
		return
	}

	likeQuery := "%" + query + "%s"
	offset := fetchNumericQuery(r, "offset", 0, math.MaxInt, 0)
	limit := fetchNumericQuery(r, "limit", 1, 1000, 50)

	// Fetch the tracks
	var tracks []SearchFoundTrack
	result := persistance.Db.Raw(`
		SELECT sp.id    as playlist_id,
			   sp.name  as playlist_name,
			   spt.track_id,
			   spt.name as track_name,
			   spt.artists,
			   spt.album,
			   spt.added_at,
			   spt.deleted_at
		FROM spotify_playlist_tracks spt
				 JOIN spotify_playlists sp on spt.spotify_playlist_id = sp.id
				 JOIN tool_user_playlists tup on sp.id = tup.spotify_playlist_id and tup.tool_user_id = 2
		WHERE spt.name LIKE ?
		   OR spt.artists LIKE ?
		   OR spt.album LIKE ?
		   OR spt.track_id = ?
		ORDER BY spt.updated_at DESC
		LIMIT ?, ?
	`, likeQuery, likeQuery, likeQuery, query, offset, limit).Scan(&tracks)

	if result.Error != nil {
		log.Println("Failed to search tracks", result.Error)
		sentry.CaptureException(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		handlers.OutputJson(w, tracks)
	}
}
