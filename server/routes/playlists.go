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
	router.HandleFunc("/playlists/combined-changelog", handlers.AuthUser(getCombinedChangelog))
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

func getCombinedChangelog(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	// Get all the playlists the user has access to
	// Find the playlist & check if the user has them in their library
	var playlists []models.ToolUserPlaylist
	result := persistance.Db.Debug().Table("tool_user_playlists").
		Where("tool_user_id = ?", user.ID).
		Find(&playlists)
	if result.Error != nil {
		log.Println(result.Error)
		writeError(w, "An error occurred with the database", http.StatusInternalServerError)
		return
	}

	// Check if we have any results
	if len(playlists) == 0 {
		handlers.OutputJson(w, []trackEvent{})
		return
	}

	// Build the where query
	var playlistIds []interface{}
	playlistWhere := ""
	for _, playlist := range playlists {
		playlistId := playlist.SpotifyPlaylistID
		playlistIds = append(playlistIds, playlistId)
		if playlistWhere != "" {
			playlistWhere += ", "
		}
		playlistWhere += "?"
	}
	playlistWhere = "spotify_playlist_id IN (" + playlistWhere + ")"

	// Fetch the changes
	offset := fetchNumericQuery(r, "offset", 0, math.MaxInt, 0)
	limit := fetchNumericQuery(r, "limit", 1, 1000, 100)

	// Check if the user has access to the playlist
	var tracks []trackEvent
	db := persistance.Db
	err := db.Debug().Raw("(? UNION ?) ORDER BY `timeline` DESC LIMIT ?, ?",
		db.Unscoped().Table("spotify_playlist_tracks").Select("*", "'added' as `type`", "added_at as `timeline`").
			Where(playlistWhere, playlistIds...),
		db.Unscoped().Table("spotify_playlist_tracks").Select("*", "'removed' as `type`", "deleted_at as `timeline`").
			Where(playlistWhere, playlistIds...).
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
