package processing

import (
	"SpotifyTool/persistance/models"
	"SpotifyTool/server/state"
	"github.com/zmb3/spotify/v2"
	"log"
	"time"
)

const (
	maxPlaylistsPerPage           = 50
	maxTracksPerPage              = 100
	fetchPlaylistsInterval        = 15 * time.Minute
	fetchSpotifyPlaylistsInterval = 2 * time.Hour
)

func HandleTasks() {
	authenticator := state.GetSpotifyAuthenticator()

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
			log.Println("Checking playlists changes for " + toolToken.ToolUser.DisplayName)
			err := doCheckPlaylistChanges(toolToken.ToolUser, client)
			if err != nil {
				log.Println(err)
				// TODO SEntry
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
	var fetchTracksForPlaylists []updatePlaylist
	newPlaylists := map[string]spotify.SimplePlaylist{}

	// Fetch all playlist from the API
	offset := 0
	for {
		// Fetch the page
		playlistsPage, err := client.GetPlaylistsForUser(ctx, user.SpotifyId, spotify.Limit(maxPlaylistsPerPage), spotify.Offset(offset))
		if err != nil {
			log.Println(err)
			// TODO: Better error handling?
			return err
		}

		// Loop through the playlists
		for _, foundPlaylist := range playlistsPage.Playlists {
			playlistId := foundPlaylist.ID.String()

			// Check if we have the playlist already
			if currentPlaylistSnapshot := currentPlaylistIdToSnapshot[playlistId]; currentPlaylistSnapshot == "" {
				// The user doesn't follow the playlist yet, this doesn't mean we don't track it yet tho. We'll need to check the database.
				newPlaylists[playlistId] = foundPlaylist
			} else {
				// The user already follows the playlist
				localPlaylist := currentPlaylistIdToPlaylist[playlistId]
				delete(currentPlaylistIdToSnapshot, playlistId)

				// Check if the playlist has changed
				if currentPlaylistSnapshot != foundPlaylist.SnapshotID {
					// Check if the playlist hasn't already been checked recently
					if isRecentlyChecked(localPlaylist) {
						// Already checked in the last [interval] minutes (probably by a different
						continue
					}

					// Snapshot has changed -> changes to the playlist have been made
					fetchTracksForPlaylists = append(fetchTracksForPlaylists, updatePlaylist{
						Local:  localPlaylist,
						Remote: foundPlaylist,
					})

					// Update the playlist meta
					localPlaylist.FromSimpleApiPlaylist(&foundPlaylist, false)
				} else {
					// Update the last checked time
					localPlaylist.SetLastCheckedToNow()
				}
				db.Save(&localPlaylist)
			}
		}

		// Go to the next page if possible
		if playlistsPage.Next != "" {
			offset += maxPlaylistsPerPage
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
			foundNewPlaylist.FromSimpleApiPlaylist(&newPlaylist, true)
			foundNewPlaylist.SnapshotId = ""
			db.Create(&foundNewPlaylist)
			currentPlaylistIdToPlaylist[playlistId] = &foundNewPlaylist
			fetchTracksForPlaylists = append(fetchTracksForPlaylists, updatePlaylist{
				Local:  &foundNewPlaylist,
				Remote: newPlaylist,
			})
		} else {
			// We already have it... Has it changed tho? (or has it already recently been checked?)
			if foundNewPlaylist.SnapshotId != newPlaylist.SnapshotID && !isRecentlyChecked(&foundNewPlaylist) {
				fetchTracksForPlaylists = append(fetchTracksForPlaylists, updatePlaylist{
					Local:  &foundNewPlaylist,
					Remote: newPlaylist,
				})
			}
		}

		// Add the playlist to the user's playlists
		currentPlaylists = append(currentPlaylists, &foundNewPlaylist)
	}

	// Update the tracks for the given playlists
	for _, updatePlaylist := range fetchTracksForPlaylists {
		if err := updateTracksOfPlaylist(client, updatePlaylist); err != nil {
			// TODO: Sentry?
			log.Println(err)
		} else {
			// Save the update to the playlist
			updatePlaylist.Local.SnapshotId = updatePlaylist.Remote.SnapshotID
			updatePlaylist.Local.SetLastCheckedToNow()
			db.Save(&updatePlaylist.Local)
		}
	}

	// Remove connection between user & playlists that are still in 'currentPlaylistIdToSnapshot', the user doesn't follow them anymore
	// TODO: Remove connection
	// TODO: Playlist access history

	// Update the relation between the user & the playlists
	if err := db.Model(&user).Association("Playlists").Replace(currentPlaylists); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Finished checking")

	return nil
}

func isRecentlyChecked(localPlaylist *models.SpotifyPlaylist) bool {
	// localPlaylist.IsAlreadyCheckedInLast(fetchPlaylistsInterval) ||
	return localPlaylist.OwnerID == "spotify" && localPlaylist.IsAlreadyCheckedInLast(fetchSpotifyPlaylistsInterval)
}

func updateTracksOfPlaylist(client *spotify.Client, update updatePlaylist) error {
	// Load the tracks of the given playlist
	if err := db.Model(update.Local).Association("Tracks").Find(&update.Local.Tracks); err != nil {
		log.Println(err)
		return err
	}

	localTracks := update.Local.Tracks
	idToLocalTrack := map[string]*models.SpotifyPlaylistTrack{}
	for _, localTrack := range localTracks {
		idToLocalTrack[localTrack.TrackId] = localTrack
	}

	// Start fetching the remote tracks
	log.Println("Updating tracks of " + update.Local.ID + ": " + update.Local.Name)

	// Fetch all playlist from the API
	tries := 0
	offset := 0
	for {
		// Fetch the page
		tracksPage, err := client.GetPlaylistTracks(ctx, update.Remote.ID,
			spotify.Limit(maxTracksPerPage), spotify.Offset(offset),
			spotify.Fields("total,next,items(added_at,added_by(id),is_local,track(id,name,artists(id,name),album(album_type,id,name)))"))
		if err != nil {
			log.Println(err)
			tries++
			if tries <= 3 {
				log.Println("=> Trying again")
				time.Sleep(5 * time.Second)
				continue
			}
			return err
		}

		// Loop through the playlists
		for _, foundTrack := range tracksPage.Tracks {
			// Skip local tracks
			if foundTrack.IsLocal || foundTrack.Track.ID.String() == "" {
				continue
			}

			if localTrack, found := idToLocalTrack[foundTrack.Track.ID.String()]; found {
				// Already exists, check if it should be updated
				if changed := localTrack.FromSpotifyPlaylistTrack(foundTrack); changed {
					db.Save(&localTrack)
				}
				delete(idToLocalTrack, localTrack.TrackId)
			} else {
				// New track
				newLocalTrack := models.SpotifyPlaylistTrack{
					SpotifyPlaylistID: update.Local.ID,
				}
				newLocalTrack.FromSpotifyPlaylistTrack(foundTrack)
				db.Create(&newLocalTrack)
				localTracks = append(localTracks, &newLocalTrack)
			}
		}

		// Go to the next page if possible
		if tracksPage.Next != "" {
			offset += maxTracksPerPage
		} else {
			break
		}
	}

	// Delete any tracks that are no longer in the playlist
	for _, deletedTrack := range idToLocalTrack {
		// Soft delete it from the DB
		db.Delete(&deletedTrack)
	}

	return nil
}

type updatePlaylist struct {
	Local  *models.SpotifyPlaylist
	Remote spotify.SimplePlaylist
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
