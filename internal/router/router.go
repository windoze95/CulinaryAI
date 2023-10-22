package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/db"
	"github.com/windoze95/saltybytes-api/internal/handlers"
	"github.com/windoze95/saltybytes-api/internal/middleware"
	"github.com/windoze95/saltybytes-api/internal/repository"
	"github.com/windoze95/saltybytes-api/internal/service"
	"golang.org/x/time/rate"
)

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
	var publicOpenAIKeyRps int = 1               // 1 request per second
	var publicOpenAIKeyBurst int = 5             // Burst of 5 requests
	var globalRps int = 20                       // 20 request per second
	var globalCleanupInterval = 10 * time.Minute // Cleanup every 10 minutes
	var globalExpiration = 1 * time.Hour         // Remove unused limiters after 1 hour

	// Define rate limiter for users with no OpenAI key
	publicOpenAIKeyRateLimiter := rate.NewLimiter(rate.Limit(publicOpenAIKeyRps), publicOpenAIKeyBurst)

	// Apply rate limiting middleware to all routes
	r.Use(middleware.RateLimitByIP(globalRps, globalCleanupInterval, globalExpiration))
	r.Use(middleware.CheckIDHeader())

	// // Individual static routes for specific files
	// r.StaticFile("/", "./web/saltybytes/build/index.html")
	// r.StaticFile("/asset-manifest.json", "./web/saltybytes/build/asset-manifest.json")
	// r.StaticFile("/favicon.ico", "./web/saltybytes/build/favicon.ico")
	// r.StaticFile("/logo192.png", "./web/saltybytes/build/logo192.png")
	// r.StaticFile("/logo512.png", "./web/saltybytes/build/logo512.png")
	// r.StaticFile("/manifest.json", "./web/saltybytes/build/manifest.json")
	// r.StaticFile("/robots.txt", "./web/saltybytes/build/robots.txt")

	// // Static route for files under "static" directory
	// r.Static("/static", "./web/saltybytes/build/static")

	// Ping route for testing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// User-related routes setup
	userDB := db.NewUserDB(database)
	userRepo := repository.NewUserRepository(userDB)
	userService := service.NewUserService(cfg, userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Recipe-related routes setup
	recipeDB := db.NewRecipeDB(database)
	recipeRepo := repository.NewRecipeRepository(recipeDB)
	recipeService := service.NewRecipeService(cfg, recipeRepo)
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Group for API routes that don't require token verification
	apiPublic := r.Group("/v1")
	{
		// Create a new user
		apiPublic.POST("/users", userHandler.CreateUser)
		// Login a user
		apiPublic.POST("/auth/login", userHandler.LoginUser)
		// Facebook OAuth routes
		apiPublic.POST("/auth/facebook", userHandler.FacebookAuth)
		apiPublic.POST("/auth/facebook/callback", userHandler.FacebookCallback)
		apiPublic.POST("/auth/facebook/complete", userHandler.CompleteFacebookSignup)
		// Get a single recipe by it's ID
		apiPublic.GET("/recipes/:recipe_id", recipeHandler.GetRecipe)
	}

	// Group for API routes that require token verification
	apiProtected := r.Group("/v1")
	{
		apiProtected.Use(middleware.VerifyTokenMiddleware(cfg))

		// User-related routes

		// Verify a user's token
		apiProtected.GET("/users/verify", middleware.AttachUserToContext(userService), userHandler.VerifyToken)
		// Logout a user
		apiProtected.POST("/users/logout", middleware.AttachUserToContext(userService), userHandler.LogoutUser)
		// Get a user by their ID
		apiProtected.GET("/users/me", middleware.AttachUserToContext(userService), userHandler.GetUserByID)
		// Get a user's settings
		apiProtected.GET("/users/settings", middleware.AttachUserToContext(userService), userHandler.GetUserSettings)
		// Update a user's settings
		apiProtected.PUT("/users/settings", middleware.AttachUserToContext(userService), userHandler.UpdateUserSettings)

		// Recipe-related routes

		// // Get a single recipe by it's ID
		// apiProtected.GET("/recipes/:recipe_id", recipeHandler.GetRecipe)
		// Create a new recipe
		apiProtected.POST("/recipes", middleware.AttachUserToContext(userService), middleware.RateLimitPublicOpenAIKey(publicOpenAIKeyRateLimiter), recipeHandler.CreateRecipe)
	}

	// // Catch-all route for serving back the React app
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("./web/saltybytes/build/index.html")
	// })

	return r
}
