package processing

import (
	"SpotifyTool/persistance/models"
	"SpotifyTool/processing/tasks"
	"SpotifyTool/server/state"
	"github.com/zmb3/spotify/v2"
	"log"
)

func HandleTasks() {
	authenticator := state.GetSpotifyAuthenticator()

	for {
		// Retrieve the first task in queue
		task := <-onSpotifyTaskChannel

		// Check if the user has a valid access-token
		toolToken := models.ToolSpotifyAuthToken{}
		db.Where("tool_user_id = ?", task.ToolUserID).Preload("ToolUser").Find(&toolToken)
		if toolToken.ID == 0 {
			// TODO: Sentry?
			// TODO: Maybe want to mail the user if the key has expired or something
			log.Println("[WARNING] Didn't find active Auth token for user: ", task.ToolUserID)
			continue
		}

		// Create the spotify client
		authToken := toolToken.ToOAuthToken()
		client := spotify.New(authenticator.Client(ctx, &authToken), spotify.WithRetry(true))

		// Do the task
		if task.Task == CheckPlaylistChanges {
			log.Println("Checking playlists changes for " + toolToken.ToolUser.DisplayName)
			err := tasks.DoCheckPlaylistChanges(toolToken.ToolUser, client)
			if err != nil {
				log.Println(err)
				// TODO SEntry
				return
			}
		} else if task.Task == BackupPlaylists {
			panic("BackupPlaylists is not supported yet")
		} else {
			panic("Impossible task given")
		}

		// Update the access token (if it has changed)
		newAuthToken, _ := client.Token()
		toolToken.FillFromOAuthToken(newAuthToken)
		db.Save(&toolToken)
	}
}

type SpotifyFetchTask struct {
	ToolUserID uint
	Task       SpotifyFetchTaskType
}

// SpotifyFetchTaskType are the types of tasks that can be executed
type SpotifyFetchTaskType int

const (
	// CheckPlaylistChanges will check all playlists & register any changes it has found since the last sync
	CheckPlaylistChanges SpotifyFetchTaskType = iota

	// BackupPlaylists will copy all the tracks of configured playlists to other specified playlists (useful for dynamic Spotify playlists)
	BackupPlaylists
)
