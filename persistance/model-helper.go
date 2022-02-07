package persistance

import "SpotifyTool/persistance/models"

func GetToolUserBySpotifyId(spotifyId string) models.ToolUser {
	user := models.ToolUser{}
	Db.Where("spotify_id = ?", spotifyId).Find(&user)

	return user
}
