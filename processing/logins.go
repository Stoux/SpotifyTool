package processing

// HandleLogins listen to logins that have taken place and inserts them into the database
func HandleLogins() {
	for {
		login := <-onLoginChannel
		go handleLogin(login)
	}
}

func handleLogin(login SpotifyClientLogin) {
	// TODO: Update most recent auth tokens

	// TODO: Check if the user is new
	// => Create profile
	// => Trigger instant full profile index
}
