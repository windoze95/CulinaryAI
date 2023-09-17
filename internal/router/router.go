package router

import (
	"time"

	"github.com/windoze95/culinaryai/internal/config"
	"github.com/windoze95/culinaryai/internal/db"
	"github.com/windoze95/culinaryai/internal/handlers"
	"github.com/windoze95/culinaryai/internal/middleware"
	"github.com/windoze95/culinaryai/internal/repository"
	"github.com/windoze95/culinaryai/internal/service"
	"golang.org/x/time/rate"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, database *gorm.DB) *gin.Engine {
	// Set Gin mode to release
	gin.SetMode(gin.ReleaseMode)

	// Create default Gin router
	r := gin.Default()

	// Define constants and variables related to rate limiting
	var publicOpenAIKeyRps int = 1               // 1 request per second
	var publicOpenAIKeyBurst int = 5             // Burst of 5 requests
	var globalRps int = 20                       // 20 request per second
	var globalCleanupInterval = 10 * time.Minute // Cleanup every 10 minutes
	var globalExpiration = 1 * time.Hour         // Remove unused limiters after 1 hour

	// Define rate limiter for users with no OpenAI key
	publicOpenAIKeyRateLimiter := rate.NewLimiter(rate.Limit(publicOpenAIKeyRps), publicOpenAIKeyBurst)

	// Apply rate limiting middleware to all routes
	r.Use(middleware.RateLimitByIP(globalRps, globalCleanupInterval, globalExpiration))

	// Serve static files
	r.Static("/", "./web/culinaryai/build/")

	// User-related routes setup
	userDB := db.NewUserDB(database)
	userRepo := repository.NewUserRepository(userDB)
	userService := service.NewUserService(cfg, userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new user
	r.POST("/api/v1/users", userHandler.CreateUser)

	// Login a user
	r.POST("/api/v1/users/login", userHandler.LoginUser)

	r.Use(middleware.VerifyTokenMiddleware(cfg))

	// Get a user's settings
	r.GET("/api/v1/users/settings", middleware.AttachUserToContext(userService), userHandler.GetUserSettings)

	// Update a user's settings
	r.PUT("/api/v1/users/settings", middleware.AttachUserToContext(userService), userHandler.UpdateUserSettings)

	// Recipe-related routes setup
	recipeDB := db.NewRecipeDB(database)
	recipeRepo := repository.NewRecipeRepository(recipeDB)
	recipeService := service.NewRecipeService(cfg, recipeRepo)
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Get a single recipe by it's ID
	r.GET("/api/v1/recipes/:recipe_id", recipeHandler.GetRecipe)

	// Create a new recipe
	r.POST("/api/v1/recipes", middleware.RateLimitPublicOpenAIKey(publicOpenAIKeyRateLimiter), middleware.AttachUserToContext(userService), recipeHandler.CreateRecipe)

	return r
}
