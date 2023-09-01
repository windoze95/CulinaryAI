package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type GlobalConfig struct {
	Title     string `json:"title"`
	Env       Env    `json:"env"`
	MfaIssuer string `json:"mfaIssuer"`
}

type Env struct {
	Port                   string `json:"port"`
	DatabaseUrl            string `json:"databaseUrl"`
	OpenAIKeyEncryptionKey string `json:"openAIKeyEncryptionKey"`
	SessionKey             string `json:"sessionKey"`
	PublicOpenAIKey        string `json:"publicOpenAIKey"`
}

var (
	gc    *GlobalConfig
	db    *gorm.DB
	store *sessions.CookieStore
	// globalLimiter ratelimit.Limiter
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[GIN] ")
	log.SetOutput(gin.DefaultWriter)

	ConfigRuntime()

	var err error

	// Load global config
	gc, err = GetGlobalConfig()
	if err != nil {
		log.Fatalf("Error loading global config: %v", err)
	}

	// Load new store
	store = sessions.NewCookieStore([]byte(os.Getenv(gc.Env.SessionKey)))

	// Check that envs are populated
	err = CheckEnvironmentVariables(&gc.Env)
	if err != nil {
		log.Fatalf("Environment variable error: %v", err)
	}

	// Check that env are valid
	encryptedKey, err := encryptOpenAIKey(os.Getenv(gc.Env.PublicOpenAIKey))
	if err != nil {
		log.Fatalf("Unable to encrypt public openai key: %v", err)
	}
	isValid, err := verifyOpenAIKey(encryptedKey)
	if err != nil {
		log.Fatalf("Error during public openai key verification: %v", err)
	}
	if !isValid {
		log.Fatalf("Invalid public openai key: %v", err)
	}

	// Connect to the database
	db, err = ConnectToDatabaseWithRetry(os.Getenv(gc.Env.DatabaseUrl))
	if err != nil {
		log.Fatalf("Environment variable error: %v", err)
	}
	db.AutoMigrate(&User{}, &UserSettings{}, &Recipe{}, &GuidingContent{}, &Tag{})
}

func main() {
	startGin()
}

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func GetGlobalConfig() (*GlobalConfig, error) {
	var config GlobalConfig
	configFile, err := os.ReadFile("global_config.json")
	if err != nil {
		return &config, err
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}

func ConnectToDatabaseWithRetry(dbURL string) (*gorm.DB, error) {
	var database *gorm.DB
	var err error

	start := time.Now()
	for {
		database, err = gorm.Open("postgres", dbURL)
		if err == nil {
			break
		}
		if time.Since(start) > 10*time.Minute {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		log.Printf("Could not connect to database, retrying...")
		time.Sleep(5 * time.Second)
	}

	// Set a 5-second timeout for all queries in this session
	db.Exec("SET statement_timeout = 5000")

	return database, err
}

func CheckEnvironmentVariables(env *Env) error {
	v := reflect.ValueOf(env).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			return fmt.Errorf("$%s must be set", v.Type().Field(i).Name)
		}
	}

	return nil
}
