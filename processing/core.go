package processing

var (
	onLoginChannel       chan SpotifyClientLogin
	onSpotifyTaskChannel chan SpotifyFetchTask
)

// Init the processing tasks that keeps track of changes & handles new client logins
func Init() {
	onLoginChannel = make(chan SpotifyClientLogin, 10000)
	onSpotifyTaskChannel = make(chan SpotifyFetchTask, 10000)

	go HandleLogins()
	go HandleTasks()
}

// GetLoginChannel fetches the channel on which a new SpotifyClientLogin can be posted
func GetLoginChannel() chan<- SpotifyClientLogin {
	return onLoginChannel
}

// GetTaskChannel returns the (buffered) channel on which new SpotifyFetchTask items can be posted
func GetTaskChannel() chan<- SpotifyFetchTask {
	return onSpotifyTaskChannel
}
