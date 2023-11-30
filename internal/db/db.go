package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/models"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	return connectToDatabaseWithRetry(cfg.Env.DatabaseUrl.Value())
}

func connectToDatabaseWithRetry(dbURL string) (*gorm.DB, error) {
	var database *gorm.DB
	var err error

	start := time.Now()
	for {
		database, err = gorm.Open("postgres", dbURL)
		if err == nil {
			break
		}
		if time.Since(start) > 10*time.Minute {
			return nil, fmt.Errorf("could not connect to database after 10 minutes: %w", err)
		}
		log.Printf("Could not connect to database, retrying...")
		time.Sleep(5 * time.Second)
	}

	// Set a 5-second timeout for all queries in this session
	// db.Exec("SET statement_timeout = 5000")

	database.AutoMigrate(
		&models.User{},
		&models.UserAuth{},
		&models.Subscription{},
		&models.UserSettings{},
		&models.GuidingContent{},
		&models.Recipe{},
		&models.Tag{},
		&models.RecipeChatHistory{},
		&models.RecipeChatHistoryMessage{},
	)

	return database, err
}
