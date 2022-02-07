package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/handlers"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify/v2"
	"log"
	"net/http"
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

	// TODO: Offset & limit
	//offset := r.URL.Query().Get("offset")
	//limit := r.URL.Query().Get("limit")

	// Check if the user has access to the playlist
	// TODO: Order by last relevant modification
	var tracks *[]models.SpotifyPlaylistTrack
	find := persistance.Db.Unscoped().
		Where("spotify_playlist_id = ?", playlistId).
		Limit(100).
		Offset(0).
		Order("added_at desc").
		Find(&tracks)

	if find.Error != nil {
		log.Println(find.Error)
	}

	handlers.OutputJson(w, tracks)
}

type trackResult struct {
	tracks *[]spotify.PlaylistTrack
}
