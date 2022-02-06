package processing

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"database/sql"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// HandleLogins listen to logins that have taken place and inserts them into the database
func HandleLogins() {
	for {
		login := <-onLoginChannel
		go handleLogin(login)
	}
}

func handleLogin(login SpotifyClientLogin) {
	db := persistance.GetDatabase()

	// Check if the user already exists
	newUser := false
	user := models.ToolUser{}
	db.Where("spotify_id = ?", login.User.ID).Find(&user)
	if user.ID == 0 {
		// New user
		user := models.ToolUser{
			SpotifyId:         login.User.ID,
			DisplayName:       login.User.DisplayName,
			Email:             nullString(login.User.Email),
			SpotifyUri:        nullString(string(login.User.URI)),
			SpotifyProfileUrl: nullString(login.User.Endpoint),
			Plan:              nullString(login.User.Product),
		}
		db.Create(&user)
		newUser = true
	}

	// Update the auth token
	token := models.ToolSpotifyAuthToken{}
	db.Where("tool_user_id = ?", user.ID).Find(&token)
	token.FillFromOAuthToken(login.Token)
	if token.ID == 0 {
		token.ToolUserID = user.ID
		db.Create(&token)
	} else {
		db.Save(&token)
	}

	// Trigger full playlist fetch if the user is new
	if newUser {
		GetTaskChannel() <- SpotifyFetchTask{
			ToolUserID: user.ID,
			Task:       CheckPlaylistChanges,
		}
	}
}

func nullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

type SpotifyClientLogin struct {
	Token *oauth2.Token
	User  *spotify.PrivateUser
}
