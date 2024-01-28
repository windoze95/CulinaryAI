package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/db"
	"github.com/windoze95/saltybytes-api/internal/router"
)

// init is called before the main function.
func init() {
	// Configure the logger
	ConfigureLogger()

	// Configure the runtime
	ConfigureRuntime()
}

// Entry point for the API.
func main() {
	// Load the config
	var cfg *config.Config
	if c, err := config.LoadConfig("configs/config.json"); err != nil {
		log.Fatalf("Error loading config: %v", err)
	} else {
		cfg = c
	}

	// Check that all ENV variables are set
	if err := cfg.CheckConfigEnvFields(); err != nil {
		log.Fatalf("Error checking config fields: %v", err)
	}

	// Load API keys and prompts
	if err := cfg.LoadOpenaiKeys(); err != nil {
		log.Fatalf("Error loading OpenAI keys: %v", err)
	}
	log.Printf("Loaded OpenAI keys: %v", cfg.OpenaiKeys)
	if err := cfg.LoadOpenaiPrompts(); err != nil {
		log.Fatalf("Error loading OpenAI prompts: %v", err)
	}

	// Connect to the database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer database.Close()

	// Create a new gin router
	r := router.SetupRouter(cfg, database)

	// Run the server
	r.Run(":" + cfg.Env.Port.Value())
}

// ConfigureLogger sets up the logging environment.
func ConfigureLogger() {
	log.SetFlags(0)
	log.SetPrefix("[GIN] ")
	log.SetOutput(gin.DefaultWriter)
}

// ConfigureRuntime sets the number of operating system threads.
func ConfigureRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}
