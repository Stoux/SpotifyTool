package models

import (
	"database/sql"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
	"log"
	"time"
)

type SpotifyPlaylist struct {
	ID               string `gorm:"primarykey;type:varchar(255)" json:"id"`
	SnapshotId       string `gorm:"type:varchar(255)"`
	Name             string `gorm:"type:varchar(255)" json:"name"`
	Public           bool   `json:"public"`
	Collaborative    bool   `json:"collaborative"`
	OwnerDisplayName string `gorm:"type:varchar(255)" json:"owner_display_name"`
	OwnerID          string `gorm:"type:varchar(255)" json:"owner_id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	LastChecked      time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Users contain the list of users that have this playlist in their library
	Users  []*ToolUser `gorm:"many2many:tool_user_playlists"`
	Tracks []*SpotifyPlaylistTrack
}

type ToolUserPlaylist struct {
	ToolUserID        int
	SpotifyPlaylistID string
}

func (sp *SpotifyPlaylist) FromSimpleApiPlaylist(playlist *spotify.SimplePlaylist, changeSnapshot bool) {
	sp.ID = playlist.ID.String()
	if changeSnapshot {
		sp.SnapshotId = playlist.SnapshotID
	}
	sp.Name = playlist.Name
	sp.Public = playlist.IsPublic
	sp.Collaborative = playlist.Collaborative
	sp.OwnerDisplayName = playlist.Owner.DisplayName
	sp.OwnerID = playlist.Owner.ID
}

func (sp *SpotifyPlaylist) SetLastCheckedToNow() {
	sp.LastChecked = time.Now()
}

// IsAlreadyCheckedInLast checks if this playlist was already checked in the last [given time duration].
// It also adds a small buffer of a couple of minutes to be a bit more lenient.
func (sp *SpotifyPlaylist) IsAlreadyCheckedInLast(d time.Duration) bool {
	const buffer = 5 * time.Minute
	if d > buffer {
		d = d - buffer
	}

	return sp.LastChecked.Add(d).After(time.Now())
}

type SpotifyPlaylistTrack struct {
	gorm.Model

	SpotifyPlaylistID string
	SpotifyPlaylist   SpotifyPlaylist

	TrackId string `gorm:"type:varchar(255)"`
	Name    string `gorm:"type:varchar(255)"`
	// Artists = pipe | seperated artist names
	Artists string `gorm:"type:varchar(1000)"`
	// Album = name of the album it's on
	Album string `gorm:"type:varchar(255)"`

	// AddedAt = Time when the track got added to the playlist
	// Might be null on very old playlists
	AddedAt sql.NullTime
	// AddedBy = Spotify User ID.
	// Might be null on very old playlists
	AddedBy sql.NullString `gorm:"type:varchar(255)"`
}

func (t *SpotifyPlaylistTrack) FromSpotifyPlaylistTrack(spt spotify.PlaylistTrack) (changed bool) {
	t.TrackId, changed = getChanged(t.TrackId, spt.Track.ID.String(), changed)
	t.Name, changed = getChanged(t.Name, spt.Track.Name, changed)
	t.Artists, changed = getChanged(t.Artists, GetCombinedArtists(&spt), changed)
	t.Album, changed = getChanged(t.Album, spt.Track.Album.Name, changed)
	if nullableAddedBy := AsNullableString(spt.AddedBy.ID); nullableAddedBy.String != t.AddedBy.String || nullableAddedBy.Valid != t.AddedBy.Valid {
		changed = true
		t.AddedBy = nullableAddedBy
	}
	if nullableAddedAt := SpotifyDateToNullableTime(spt.AddedAt); nullableAddedAt.Time != t.AddedAt.Time || nullableAddedAt.Valid != t.AddedAt.Valid {
		changed = true
		t.AddedAt = nullableAddedAt
	}
	return changed
}

func SpotifyDateToNullableTime(sptTime string) sql.NullTime {
	if sptTime != "" && sptTime != "1970-01-01T00:00:00Z" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z", sptTime); err == nil {
			return sql.NullTime{
				Time:  parsed,
				Valid: true,
			}
		} else {
			log.Println(err)
		}
	}

	return sql.NullTime{}
}

func GetCombinedArtists(track *spotify.PlaylistTrack) string {
	result := ""
	for _, artist := range track.Track.Artists {
		if result != "" {
			result += " | "
		}
		result += artist.Name
	}
	return result
}

func getChanged(original string, new string, changed bool) (string, bool) {
	if original != new {
		return new, true
	} else {
		return original, changed
	}
}
