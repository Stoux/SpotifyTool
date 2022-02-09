package models

import (
	"gorm.io/gorm"
	"time"
)

type PlaylistBackupConfig struct {
	gorm.Model

	ToolUserId uint `gorm:"uniqueIndex:idx_unique_config"`
	ToolUser   ToolUser

	SourcePlaylistID string          `gorm:"type:varchar(191);uniqueIndex:idx_unique_config"`
	SourcePlaylist   SpotifyPlaylist `json:"source"`

	TargetPlaylistID string          `gorm:"type:varchar(191);uniqueIndex:idx_unique_config"`
	TargetPlaylist   SpotifyPlaylist `json:"target"`

	LastSync time.Time `json:"last_sync"`
	Comment  string    `grom:"type:varchar(1000)" json:"comment"`
}
