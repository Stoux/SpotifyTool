package processing

import (
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

var (
	onLoginChannel chan SpotifyClientLogin
)

// Init the processing tasks that keeps track of changes & handles new client logins
func Init() {
	onLoginChannel = make(chan SpotifyClientLogin, 10000)

	go HandleLogins()
}

// GetLoginChannel fetches the channel on which a new SpotifyClientLogin can be posted
func GetLoginChannel() chan SpotifyClientLogin {
	return onLoginChannel
}

type SpotifyClientLogin struct {
	token *oauth2.Token
	user  spotify.User
}
