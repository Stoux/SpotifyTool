package processing

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/processing/constants"
	"SpotifyTool/server/state"
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	onLoginChannel       chan SpotifyClientLogin
	onSpotifyTaskChannel chan SpotifyFetchTask
	ctx                  = context.Background()
	db                   *gorm.DB
)

// Init the processing tasks that keeps track of changes & handles new client logins
func Init() {
	onLoginChannel = make(chan SpotifyClientLogin, 10000)
	onSpotifyTaskChannel = make(chan SpotifyFetchTask, 10000)
	db = persistance.Db

	go HandleLogins()
	go HandleTasks()

	startRecurringTasks()
}

// GetLoginChannel fetches the channel on which a new SpotifyClientLogin can be posted
func GetLoginChannel() chan<- SpotifyClientLogin {
	return onLoginChannel
}

// GetTaskChannel returns the (buffered) channel on which new SpotifyFetchTask items can be posted
func GetTaskChannel() chan<- SpotifyFetchTask {
	return onSpotifyTaskChannel
}

func startRecurringTasks() {
	if hasRecurring := state.GetBoolLikeEnv("SPOTIFY_TOOL_START_RECURRING_TASKS", true); !hasRecurring {
		log.Println("Recurring tasks are disabled through due to ENV SPOTIFY_TOOL_START_RECURRING_TASKS")
		return
	}

	ticker := time.NewTicker(constants.FetchPlaylistsInterval)
	go func() {
		for ; true; <-ticker.C {
			// Fetch all tokens
			var tokens []*models.ToolSpotifyAuthToken
			db.Find(&tokens)
			log.Println("[TICK] Playlist interval")

			// Schedule fetch tasks for them
			for _, token := range tokens {
				log.Println("Triggering task for:", token.ToolUserID)
				onSpotifyTaskChannel <- SpotifyFetchTask{
					ToolUserID: token.ToolUserID,
					Task:       CheckPlaylistChanges,
				}
			}
		}
	}()
}
