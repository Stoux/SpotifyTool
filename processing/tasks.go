package processing

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/state"
	"context"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
	"log"
)

var (
	ctx = context.Background()
	db  *gorm.DB
)

const maxPerPage = 50

func HandleTasks() {
	authenticator := state.GetSpotifyAuthenticator()
	db = persistance.GetDatabase()

	for {
		// Retrieve the first task in queue
		task := <-onSpotifyTaskChannel

		// Check if the user has a valid access-token
		toolToken := models.ToolSpotifyAuthToken{}
		db.Where("tool_user_id = ?", task.ToolUserID).Preload("ToolUser").Find(&toolToken)
		if toolToken.ID == 0 {
			// TODO: Sentry?
			// TODO: Maybe want to mail the user if the key has expired or something
			log.Println("[WARNING] Didn't find active Auth token for user: ", task.ToolUserID)
			continue
		}

		// Create the spotify client
		authToken := toolToken.ToOAuthToken()
		client := spotify.New(authenticator.Client(ctx, &authToken), spotify.WithRetry(true))

		// Do the task
		if task.Task == CheckPlaylistChanges {
			err := doCheckPlaylistChanges(toolToken.ToolUser, client)
			if err != nil {
				// TODO
				return
			}
		} else if task.Task == BackupPlaylists {
			panic("BackupPlaylists is not supported yet")
		} else {
			panic("Impossible task given")
		}

		// Update the access token (if it has changed)
		newAuthToken, _ := client.Token()
		toolToken.FillFromOAuthToken(newAuthToken)
		db.Save(&toolToken)
	}
}

func doCheckPlaylistChanges(user models.ToolUser, client *spotify.Client) (err error) {
	// Fetch all playlists this user already has in our database
	db.Model(&user).Association("Playlists").Find(&user.Playlists)
	currentPlaylists := user.Playlists

	// Map to their snapshot ID
	currentPlaylistIdToPlaylist := map[string]*models.SpotifyPlaylist{}
	currentPlaylistIdToSnapshot := map[string]string{}
	for _, currentPlaylist := range currentPlaylists {
		currentPlaylistIdToPlaylist[currentPlaylist.ID] = currentPlaylist
		currentPlaylistIdToSnapshot[currentPlaylist.ID] = currentPlaylist.SnapshotId
	}

	// Build a map of playlists that we have to update as their snapshot ID has changed
	var fetchTracksForPlaylists []spotify.SimplePlaylist
	newPlaylists := map[string]spotify.SimplePlaylist{}

	// Fetch all playlist from the API
	offset := 0
	for {
		// Fetch the page
		playlistsPage, err := client.GetPlaylistsForUser(ctx, user.SpotifyId, spotify.Limit(maxPerPage), spotify.Offset(offset))
		if err != nil {
			return err
		}

		// Loop through the playlists
		for _, foundPlaylist := range playlistsPage.Playlists {
			playlistId := foundPlaylist.ID.String()

			// Check if we have the playlist already
			if currentPlaylistSnapshot := currentPlaylistIdToSnapshot[playlistId]; currentPlaylistSnapshot == "" {
				// The user doesn't follow the playlist yet, this doesn't mean we don't track it yet tho.
				newPlaylists[playlistId] = foundPlaylist
			} else if currentPlaylistSnapshot != foundPlaylist.SnapshotID {
				// Snapshot has changed -> changes to the playlist have been made
				fetchTracksForPlaylists = append(fetchTracksForPlaylists, foundPlaylist)
				delete(currentPlaylistIdToSnapshot, playlistId)

				// Update the playlist meta
				currentPlaylist := currentPlaylistIdToPlaylist[playlistId]
				currentPlaylist.FromSimpleApiPlaylist(&foundPlaylist)
				db.Save(&currentPlaylist) // Risky to update the snapshot before checking the changes?
			}
		}

		// Go to the next page if possible
		if playlistsPage.Next != "" {
			offset += maxPerPage
		} else {
			break
		}
	}

	// Check if the new playlists are actually new or that we follow them
	for playlistId, newPlaylist := range newPlaylists {
		// Fetch the playlist from our database
		foundNewPlaylist := models.SpotifyPlaylist{ID: playlistId}
		db.Find(&foundNewPlaylist)
		if foundNewPlaylist.SnapshotId == "" {
			// Actually new playlist
			foundNewPlaylist.FromSimpleApiPlaylist(&newPlaylist)
			db.Create(&foundNewPlaylist)
		} else {
			// We already have it... Has it changed tho?
			if foundNewPlaylist.SnapshotId != newPlaylist.SnapshotID {
				fetchTracksForPlaylists = append(fetchTracksForPlaylists, newPlaylist)
			}
		}

		// Add the playlist to the user's playlists
		currentPlaylists = append(currentPlaylists, &foundNewPlaylist)
	}

	// Update the tracks for the given playlists
	//page, err := client.GetPlaylistTracks(ctx, "")
	//page.Tracks[0].

	// Remove connection between user & playlists that are still in 'currentPlaylistIdToSnapshot', the user doesn't follow them anymore
	// TODO: Remove connection
	// TODO: Playlist access history

	// Update the relation between the user & the playlists
	if err := db.Model(&user).Association("Playlists").Replace(currentPlaylists); err != nil {
		return err
	}

	// Update the access token (if it has changed)

	return nil
}

func updateTracksOfPlaylist() {

}

type SpotifyFetchTask struct {
	ToolUserID uint
	Task       SpotifyFetchTaskType
}

// SpotifyFetchTaskType are the types of tasks that can be executed
type SpotifyFetchTaskType int

const (
	// CheckPlaylistChanges will check all playlists & register any changes it has found since the last sync
	CheckPlaylistChanges SpotifyFetchTaskType = iota

	// BackupPlaylists will copy all the tracks of configured playlists to other specified playlists (useful for dynamic Spotify playlists)
	BackupPlaylists
)
