package persistance

import (
	"SpotifyTool/persistance/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	Db *gorm.DB
)

// Init the database connection & ORM.
func Init() {
	// Resolve the DB variables
	var (
		user         = getEnvOrDefault("SPOTIFY_TOOL_MYSQL_USER", "root")
		password     = getEnvOrDefault("SPOTIFY_TOOL_MYSQL_PASSWORD", "root")
		databaseName = getEnvOrDefault("SPOTIFY_TOOL_MYSQL_DATABASE", "spotify-tool")
		protocol     = getEnvOrDefault("SPOTIFY_TOOL_MYSQL_PROTOCOL", "tcp")
		address      = getEnvOrDefault("SPOTIFY_TOOL_MYSQL_ADDRESS", "127.0.0.1:3306")
	)

	// Init the DB
	var err error
	dsn := user + ":" + password + "@" + protocol + "(" + address + ")/" + databaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Create the DB columns (if not done yet) by migrating all models
	migrateSchema()
}

func getEnvOrDefault(env string, defaultValue string) string {
	if value := os.Getenv(env); value != "" {
		return value
	} else {
		return defaultValue
	}
}

func migrateSchema() {
	err := Db.AutoMigrate(
		models.ToolUser{},
		models.ToolSpotifyAuthToken{},
		models.SpotifyPlaylist{},
		models.SpotifyPlaylistTrack{},
		models.PlaylistBackupConfig{},
	)
	if err != nil {
		panic(err)
	}
}
