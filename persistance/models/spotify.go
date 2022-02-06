package models

import (
	"database/sql"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
	"time"
)

type SpotifyPlaylist struct {
	ID               string `gorm:"primarykey"`
	SnapshotId       string `gorm:"type:varchar(255)"`
	Name             string `gorm:"type:varchar(255)"`
	Public           bool
	Collaborative    bool
	OwnerDisplayName string `gorm:"type:varchar(255)"`
	OwnerID          string `gorm:"type:varchar(255)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	LastChecked      time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Users contain the list of users that have this playlist in their library
	Users []*ToolUser `gorm:"many2many:tool_user_playlists"`
}

func (sp *SpotifyPlaylist) FromSimpleApiPlaylist(playlist *spotify.SimplePlaylist) {
	sp.ID = playlist.ID.String()
	sp.SnapshotId = playlist.SnapshotID
	sp.Name = playlist.Name
	sp.Public = playlist.IsPublic
	sp.Collaborative = playlist.Collaborative
	sp.OwnerDisplayName = playlist.Owner.DisplayName
	sp.OwnerID = playlist.Owner.ID
}

type SpotifyTrackInPlaylist struct {
	gorm.Model

	AddedAt sql.NullTime
	AddedBy sql.NullString `gorm:"type:varchar(255)"`
}