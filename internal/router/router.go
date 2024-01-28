package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/handlers"
	"github.com/windoze95/saltybytes-api/internal/middleware"
	"github.com/windoze95/saltybytes-api/internal/repository"
	"github.com/windoze95/saltybytes-api/internal/service"
)

// SetupRouter sets up the Gin router.
func SetupRouter(cfg *config.Config, database *gorm.DB) *gin.Engine {
	// Set Gin mode to release
	gin.SetMode(gin.ReleaseMode)

	// Create default Gin router
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{
		"https://api.saltybytes.ai",
		"https://www.api.saltybytes.ai",
		"https://saltybytes.ai",
		"https://www.saltybytes.ai",
	}
	config.AllowHeaders = append(config.AllowHeaders, "X-SaltyBytes-Identifier")

	r.Use(cors.New(config))

	// Define constants and variables related to rate limiting
	var globalRps int = 20                       // 20 request per second
	var globalCleanupInterval = 10 * time.Minute // Cleanup every 10 minutes
	var globalExpiration = 1 * time.Hour         // Remove unused limiters after 1 hour

	// Apply rate limiting middleware to all routes
	r.Use(middleware.RateLimitByIP(globalRps, globalCleanupInterval, globalExpiration))
	r.Use(middleware.CheckIDHeader(cfg.Env.IdHeader.Value()))

	// Ping route for testing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// User-related routes setup
	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(cfg, userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Recipe-related routes setup
	recipeRepo := repository.NewRecipeRepository(database)
	recipeService := service.NewRecipeService(cfg, recipeRepo)
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Group for API routes that don't require token verification
	apiPublic := r.Group("/v1")
	{
		// User-related routes

		// Create a new user
		apiPublic.POST("/users", userHandler.CreateUser)
		// Login a user
		apiPublic.POST("/auth/login", userHandler.LoginUser)

		// Recipe-related routes

		// Get a single recipe by it's ID
		apiPublic.GET("/recipes/:recipe_id", recipeHandler.GetRecipe)
		// Get a single recipe history by the recipe history's ID
		apiPublic.GET("/recipes/chat-history/:history_id", recipeHandler.GetRecipeHistory)
	}

	// Group for API routes that require token verification
	apiProtected := r.Group("/v1")
	{
		apiProtected.Use(middleware.VerifyTokenMiddleware(cfg))

		// User-related routes

		// Verify a user's token
		apiProtected.GET("/users/verify", middleware.AttachUserToContext(userService), userHandler.VerifyToken)
		// Get a user by their ID
		apiProtected.GET("/users/me", middleware.AttachUserToContext(userService), userHandler.GetUserByID)
		// Get a user's settings
		apiProtected.GET("/users/settings", middleware.AttachUserToContext(userService), userHandler.GetUserSettings)

		// Recipe-related routes

		// // Get a single recipe by it's ID
		// apiProtected.GET("/recipes/:recipe_id", recipeHandler.GetRecipe)
		// Generate a new recipe
		apiProtected.POST("/recipes/chat", middleware.AttachUserToContext(userService), recipeHandler.GenerateRecipeWithChat)
		// Import a recipe with a link
		// apiProtected.POST("/recipes/import/link", middleware.AttachUserToContext(userService), recipeHandler.ImportRecipeLink)
		// Import a recipe with vision
		// apiProtected.POST("/recipes/import/vision", middleware.AttachUserToContext(userService), recipeHandler.ImportRecipeVision)
		// Import a recipe with copy-paste
		// apiProtected.POST("/recipes/import/copypasta", middleware.AttachUserToContext(userService), recipeHandler.ImportRecipeCopyPasta)
		// Manually enter a new recipe
		// apiProtected.POST("/recipes/manual", middleware.AttachUserToContext(userService), recipeHandler.ManualEntryRecipe)
		// Copycat a recipe
		// apiProtected.POST("/recipes/copycat", middleware.AttachUserToContext(userService), recipeHandler.CopycatRecipe)
	}

	return r
}
