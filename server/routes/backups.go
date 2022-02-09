package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/handlers"
	json2 "encoding/json"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func PlaylistBackupRoutes(router *mux.Router) {
	sub := router.PathPrefix("/playlist-backups").Subrouter()
	sub.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	sub.Methods("GET").HandlerFunc(handlers.AuthUser(getPlaylistBackups))
	sub.Methods("POST").HandlerFunc(handlers.AuthUser(createPlaylistBackup))
}

func getPlaylistBackups(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	// Load the available configs
	var configs []*models.PlaylistBackupConfig
	persistance.Db.Preload("SourcePlaylist").Preload("TargetPlaylist").
		Where("tool_user_id = ?", user.ID).
		Order("updated_at desc").
		Find(&configs)

	handlers.OutputJson(w, configs)
}

func createPlaylistBackup(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	// Read the request
	var body createBackupConfigRequest
	if err := json2.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	} else if body.Target == body.Source || body.Target == "" || body.Source == "" {
		writeError(w, "Please specify two different playlists", http.StatusBadRequest)
		return
	} else if len(body.Comment) > 1000 {
		writeError(w, "Comment cannot be longer than 1000 chars", http.StatusBadRequest)
		return
	}

	// Find the playlist & check if the user has them in their library
	var playlists []models.SpotifyPlaylist
	result := persistance.Db.Debug().Table("spotify_playlists").
		Joins("join tool_user_playlists on spotify_playlist_id = id and tool_user_id = ?", user.ID).
		Where("id = ?", body.Source).Or("id = ?", body.Target).
		Find(&playlists)
	if result.Error != nil {
		log.Println(result.Error)
		writeError(w, "An error occurred with the database", http.StatusInternalServerError)
		return
	}

	source := findPlaylist(body.Source, playlists)
	target := findPlaylist(body.Target, playlists)
	if source == nil || target == nil {
		writeError(w, "You do not have access to both playlists", http.StatusForbidden)
		return
	}

	// Check if the user has write-access to the target
	if !target.Collaborative && target.OwnerID != user.SpotifyId {
		writeError(w, "You do not have write-access to the target playlist", http.StatusForbidden)
		return
	}

	// Create the
	config := models.PlaylistBackupConfig{
		ToolUserId:       user.ID,
		SourcePlaylistID: source.ID,
		TargetPlaylistID: target.ID,
		LastSync:         time.Now(),
		Comment:          body.Comment,
	}
	if createResult := persistance.Db.Create(&config); createResult.Error != nil {
		mysqlError := &mysql.MySQLError{}
		if errors.As(createResult.Error, &mysqlError) && mysqlError.Number == 1062 {
			writeError(w, "This combination already exists, update the existing config.", http.StatusConflict)
		} else {
			log.Println(createResult.Error)
			writeError(w, "An internal database error occurred, try again later.", http.StatusInternalServerError)
		}
		return
	}

	handlers.OutputJson(w, config)
}

func findPlaylist(id string, playlists []models.SpotifyPlaylist) *models.SpotifyPlaylist {
	for _, playlist := range playlists {
		if playlist.ID == id {
			return &playlist
		}
	}
	return nil
}

func writeError(w http.ResponseWriter, error string, statusCode int) {
	w.WriteHeader(statusCode)
	handlers.OutputJson(w, map[string]string{"error": error})
}

type createBackupConfigRequest struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Comment string `json:"comment"`
}
