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
	SpotifyId         string         `gorm:"type:varchar(255);uniqueIndex" json:"spotify_id"`
	DisplayName       string         `gorm:"type:varchar(255)" json:"display_name"`
	Email             sql.NullString `gorm:"type:varchar(255)"`
	SpotifyUri        sql.NullString `gorm:"type:varchar(1000)" json:"spotify_uri"`
	SpotifyProfileUrl sql.NullString `gorm:"type:varchar(1000)" json:"spotify_profile_url"`
	// Plan is like 'premium', 'free', etc.
	Plan sql.NullString `gorm:"type:varchar(255)" json:"plan"`

	// Relations
	// Playlists contain all playlists that this user has in their library (or at least that we can see)
	Playlists             []*ToolUserPlaylist     `gorm:"constraint:OnDelete:CASCADE;"`
	PlaylistBackupConfigs []*PlaylistBackupConfig `gorm:"constraint:OnDelete:CASCADE"`
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
	if oauthToken.RefreshToken != "" {
		token.RefreshToken = oauthToken.RefreshToken
	}
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
