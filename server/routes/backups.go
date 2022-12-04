package routes

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/handlers"
	json2 "encoding/json"
	"errors"
	"github.com/getsentry/sentry-go"
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
	sub.Methods("PUT").Path("/{id:[0-9]+}").HandlerFunc(handlers.AuthUser(updatePlaylistBackup))
	sub.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(handlers.AuthUser(deletePlaylistBackup))
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
	body, source, target, success := validateBackupConfigRequest(w, r, user)
	if !success {
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
			log.Println("Failed to create playlist", createResult.Error)
			sentry.CaptureException(createResult.Error)
			writeError(w, "An internal database error occurred, try again later.", http.StatusInternalServerError)
		}
		return
	}

	handlers.OutputJson(w, config)
}

func updatePlaylistBackup(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	config, success := findConfigFromPath(w, r, user)
	if !success {
		return
	}

	body, source, target, success := validateBackupConfigRequest(w, r, user)
	if !success {
		return
	}

	config.SourcePlaylistID = source.ID
	config.TargetPlaylistID = target.ID
	config.Comment = body.Comment
	config.LastSync = time.Now()

	if saveResult := persistance.Db.Save(&config); saveResult.Error != nil {
		log.Println("Update playlist backup failed", saveResult.Error)
		sentry.CaptureException(saveResult.Error)
		writeError(w, "An internal database error occurred, try again later.", http.StatusInternalServerError)
		return
	}

	handlers.OutputJson(w, config)
}

func deletePlaylistBackup(w http.ResponseWriter, r *http.Request, user models.ToolUser) {
	config, success := findConfigFromPath(w, r, user)
	if !success {
		return
	}

	if deleteResult := persistance.Db.Delete(&config); deleteResult.Error != nil {
		log.Println("Failed to delete backup", deleteResult.Error)
		sentry.CaptureException(deleteResult.Error)
		writeError(w, "An internal database error occurred, try again later.", http.StatusInternalServerError)
	} else {
		handlers.OutputJson(w, map[string]bool{"success": true})
	}
}

func findConfigFromPath(w http.ResponseWriter, r *http.Request, user models.ToolUser) (config models.PlaylistBackupConfig, success bool) {
	vars := mux.Vars(r)
	config = models.PlaylistBackupConfig{}
	if find := persistance.Db.Where("id = ?", vars["id"]).Find(&config); find.Error != nil {
		log.Println("Failed to find config", find.Error)
		sentry.CaptureException(find.Error)
		writeError(w, "An internal database error occurred, try again later.", http.StatusInternalServerError)
		return
	}

	// Check if the user has correct rights
	if config.ToolUserId != user.ID {
		writeError(w, "You do not have access to that config.", http.StatusForbidden)
		return
	}

	return config, true
}

func validateBackupConfigRequest(w http.ResponseWriter, r *http.Request, user models.ToolUser) (
	body backupConfigRequest,
	source *models.SpotifyPlaylist,
	target *models.SpotifyPlaylist,
	success bool,
) {
	// Read the request
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
		log.Println("Failed to validate backup", result.Error)
		sentry.CaptureException(result.Error)
		writeError(w, "An error occurred with the database", http.StatusInternalServerError)
		return
	}

	source = findPlaylist(body.Source, playlists)
	target = findPlaylist(body.Target, playlists)
	if source == nil || target == nil {
		writeError(w, "You do not have access to both playlists", http.StatusForbidden)
		return
	}

	// Check if the user has write-access to the target
	if !target.Collaborative && target.OwnerID != user.SpotifyId {
		writeError(w, "You do not have write-access to the target playlist", http.StatusForbidden)
		return
	}

	return body, source, target, true
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

type backupConfigRequest struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Comment string `json:"comment"`
}
