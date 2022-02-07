package processing

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"context"
	"gorm.io/gorm"
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
	db = persistance.GetDatabase()

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
	ticker := time.NewTicker(fetchPlaylistsInterval)
	go func() {
		for ; true; <-ticker.C {
			// Fetch all tokens
			var tokens []*models.ToolSpotifyAuthToken
			db.Find(&tokens)

			// Schedule fetch tasks for them
			for _, token := range tokens {
				onSpotifyTaskChannel <- SpotifyFetchTask{
					ToolUserID: token.ToolUserID,
					Task:       CheckPlaylistChanges,
				}
			}
		}
	}()
}
