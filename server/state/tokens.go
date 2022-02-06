package state

import (
	"SpotifyTool/server/util"
	"time"
)

type AuthToken struct {
	Identifier string    `json:"identifier"`
	SpotifyId  string    `json:"spotify_id"`
	Expires    time.Time `json:"expires"`
}

var (
	authorizedTokens = map[string]AuthToken{}
	initialized      = false
)

func Init() {
	expectInit(false)

	// Run script to check for outdated auth tokens
	go func() {
		for true {
			// Remove outdated auth tokens
			now := time.Now().Unix()
			for identifier, token := range authorizedTokens {
				if token.Expires.Unix() < now {
					delete(authorizedTokens, identifier)
				}
			}

			time.Sleep(15 * time.Minute)
		}
	}()

	initialized = true
}

// CreateTokenFor will create a new AuthToken for the given Spotify ID
func CreateTokenFor(spotifyId string) *AuthToken {
	// Create the auth token
	token := AuthToken{
		Identifier: util.GetRandomString(64),
		SpotifyId:  spotifyId,
		Expires:    time.Now().Add(time.Hour * 24 * 7),
	}

	authorizedTokens[token.Identifier] = token

	return &token
}

func GetTokenBy(identifier string) *AuthToken {
	if token, found := authorizedTokens[identifier]; found {
		return &token
	} else {
		return nil
	}
}

func expectInit(init bool) {
	if init != initialized {
		if init {
			panic("Expected Users state management to be initialized but it isn't!")
		} else {
			panic("Expected Users state management to not be initialized but it already is!")
		}
	}
}
