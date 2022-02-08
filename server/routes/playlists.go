package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/handlers"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func PlaylistRoutes(router *mux.Router) {
	router.HandleFunc("/playlists", handlers.AuthUser(getPlaylists))
	router.HandleFunc("/playlists/{id}/tracks", handlers.AuthUser(getPlaylistTracks))
}

func getPlaylists(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	persistance.Db.Model(&user).Association("Playlists").Find(&user.Playlists)
	handlers.OutputJson(w, user.Playlists)
}

func getPlaylistTracks(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	// Check if the user has access to the given playlist
	playlistId := mux.Vars(r)["id"]
	var count int64
	persistance.Db.Table("tool_user_playlists").
		Where("tool_user_id = ?", user.ID).
		Where("spotify_playlist_id = ?", playlistId).
		Count(&count)
	if count < 1 {
		w.WriteHeader(403)
		return
	}

	offset := fetchNumericQuery(r, "offset", 0, math.MaxInt, 0)
	limit := fetchNumericQuery(r, "limit", 1, 1000, 100)

	// Check if the user has access to the playlist
	var tracks []trackEvent
	db := persistance.Db
	err := db.Raw("(? UNION ?) ORDER BY `timeline` DESC LIMIT ?, ?",
		db.Unscoped().Table("spotify_playlist_tracks").Select("*", "'added' as `type`", "added_at as `timeline`").
			Where("spotify_playlist_id = ?", playlistId),
		db.Unscoped().Table("spotify_playlist_tracks").Select("*", "'removed' as `type`", "deleted_at as `timeline`").
			Where("spotify_playlist_id = ?", playlistId).
			Where("deleted_at IS NOT NULL"),
		offset,
		limit,
	).Scan(&tracks)
	if err != nil {
		log.Println(err)
	}
	handlers.OutputJson(w, tracks)
}

func fetchNumericQuery(r *http.Request, key string, minValue int, maxValue int, defaultValue int) int {
	if arg := r.URL.Query().Get(key); arg != "" {
		if v, err := strconv.Atoi(r.URL.Query().Get(key)); err == nil && v >= minValue && v <= maxValue {
			return v
		}
	}

	return defaultValue
}

type trackEvent struct {
	models.SpotifyPlaylistTrack
	Type     string    `json:"type"`
	Timeline time.Time `json:"timeline"`
}
