package models

import (
	"database/sql"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

// ToolUser is a single user that is using the Spotify tool.
// Each user must have a S
type ToolUser struct {
	gorm.Model
	SpotifyId         string `gorm:"uniqueIndex"`
	DisplayName       string
	Email             sql.NullString
	SpotifyUri        sql.NullString
	SpotifyProfileUrl sql.NullString
	Plan              sql.NullString
}

// ToolSpotifyAuthToken contains OAUth tokens used to access the Spotify API
type ToolSpotifyAuthToken struct {
	gorm.Model

	ToolUserID uint
	ToolUser   ToolUser

	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func (token *ToolSpotifyAuthToken) FillFromOAuthToken(oauthToken *oauth2.Token) {
	token.AccessToken = oauthToken.AccessToken
	token.TokenType = oauthToken.TokenType
	token.RefreshToken = oauthToken.RefreshToken
	token.Expiry = oauthToken.Expiry
}

func (token *ToolSpotifyAuthToken) ToOAuthToken() oauth2.Token {
	return oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
}
