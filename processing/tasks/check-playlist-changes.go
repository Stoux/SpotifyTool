package tasks

import (
	"SpotifyTool/persistance"
	"SpotifyTool/persistance/models"
	"SpotifyTool/processing/constants"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

func DoCheckPlaylistChanges(user models.ToolUser, client *spotify.Client) (err error) {
	db = persistance.Db
	ctx = context.Background()

	// Fetch all playlists this user already has in our database
	db.Model(&user).Debug().Preload("SpotifyPlaylist").Association("Playlists").Find(&user.Playlists)
	userPlaylists := user.Playlists
	currentPlaylists := make([]models.SpotifyPlaylist, len(userPlaylists))
	isPlaylistChecked := map[string]bool{}
	for i, userPlaylist := range userPlaylists {
		currentPlaylists[i] = userPlaylist.SpotifyPlaylist
		isPlaylistChecked[userPlaylist.SpotifyPlaylistID] = userPlaylist.IsTracked
	}

	// Map to their snapshot ID
	currentPlaylistIdToPlaylist := map[string]*models.SpotifyPlaylist{}
	currentPlaylistIdToSnapshot := map[string]string{}
	for _, currentPlaylist := range currentPlaylists {
		currentPlaylistIdToPlaylist[currentPlaylist.ID] = &currentPlaylist
		currentPlaylistIdToSnapshot[currentPlaylist.ID] = currentPlaylist.SnapshotId
	}

	// Build a map of playlists that we have to update as their snapshot ID has changed
	var fetchTracksForPlaylists []updatePlaylist
	newPlaylists := map[string]spotify.SimplePlaylist{}

	// Fetch all playlist from the API
	offset := 0
	for {
		// Fetch the page
		playlistsPage, err := client.GetPlaylistsForUser(ctx, user.SpotifyId, spotify.Limit(constants.MaxPlaylistsPerPage), spotify.Offset(offset))
		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
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
					// Does the user want to check this list?
					if !isPlaylistChecked[foundPlaylist.ID.String()] {
						continue
					}

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
			offset += constants.MaxPlaylistsPerPage
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
			foundNewPlaylist.SnapshotId = "-"
			foundNewPlaylist.LastChecked = time.Unix(0, 0)
			db.Create(&foundNewPlaylist)
			if foundNewPlaylist.ShouldBeCheckedByDefault() {
				fetchTracksForPlaylists = append(fetchTracksForPlaylists, updatePlaylist{
					Local:  &foundNewPlaylist,
					Remote: newPlaylist,
				})
			}
		} else {
			// We already have it... Has it changed tho? (or has it already recently been checked?) and do we want to check it?
			if foundNewPlaylist.ShouldBeCheckedByDefault() && foundNewPlaylist.SnapshotId != newPlaylist.SnapshotID && !isRecentlyChecked(&foundNewPlaylist) {
				fetchTracksForPlaylists = append(fetchTracksForPlaylists, updatePlaylist{
					Local:  &foundNewPlaylist,
					Remote: newPlaylist,
				})
			}
		}

		// Add the playlist to the user's playlists
		isPlaylistChecked[foundNewPlaylist.ID] = foundNewPlaylist.ShouldBeCheckedByDefault()
		currentPlaylists = append(currentPlaylists, foundNewPlaylist)
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
	if len(currentPlaylistIdToSnapshot) > 0 {
		newPosition := 0
		for _, playlist := range currentPlaylists {
			if _, deleted := currentPlaylistIdToSnapshot[playlist.ID]; !deleted {
				currentPlaylists[newPosition] = playlist
				newPosition++
			}
		}
		currentPlaylists = currentPlaylists[:newPosition]
	}

	// TODO: Playlist access history

	err2 := updateUserPlaylists(user, currentPlaylists, isPlaylistChecked)
	if err2 != nil {
		sentry.CaptureException(err2)
		return err2
	}

	log.Println("Finished checking")

	return nil
}

func updateUserPlaylists(user models.ToolUser, currentPlaylists []models.SpotifyPlaylist, isPlaylistChecked map[string]bool) error {
	// Update the relation between the user & the playlists
	// => Map back to UserPlaylist entries
	newUserPlaylists := make([]models.ToolUserPlaylist, len(currentPlaylists))
	playlistIds := make([]string, len(newUserPlaylists))
	for i, newPlaylist := range currentPlaylists {
		newUserPlaylists[i] = models.ToolUserPlaylist{
			ToolUserID:        user.ID,
			SpotifyPlaylistID: newPlaylist.ID,
			IsTracked:         isPlaylistChecked[newPlaylist.ID],
		}
		playlistIds[i] = newPlaylist.ID
	}

	// Delete any playlists no longer tracked
	if result := db.Debug().Where("tool_user_id = ?", user.ID).Where("spotify_playlist_id NOT IN ?", playlistIds).Delete(&models.ToolUserPlaylist{}); result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tool_user_id"}, {Name: "spotify_playlist_id"}},
		DoUpdates: nil,
	}).Create(newUserPlaylists).Error
}

func isRecentlyChecked(localPlaylist *models.SpotifyPlaylist) bool {
	// localPlaylist.IsAlreadyCheckedInLast(fetchPlaylistsInterval) ||
	return localPlaylist.OwnerID == "spotify" && localPlaylist.IsAlreadyCheckedInLast(constants.FetchSpotifyOwnedPlaylistsInterval)
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
			spotify.Limit(constants.MaxTracksPerPage), spotify.Offset(offset),
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
			offset += constants.MaxTracksPerPage
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
