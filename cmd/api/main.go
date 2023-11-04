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

func init() {
	ConfigureLogger()
	ConfigureRuntime()
}

func main() {
	cfg, err := config.LoadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = config.CheckConfigFields(cfg)
	if err != nil {
		log.Fatalf("Error checking config fields: %v", err)
	}

	// Connect to the database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer database.Close()

	// // Load new store
	// store := sessions.NewCookieStore([]byte(cfg.Env.SessionKey.Value()))
	// store.Options = &sessions.Options{
	// 	Path: "/",
	// 	// MaxAge:   86400 * 7, // 7 days
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteStrictMode,
	// }

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
