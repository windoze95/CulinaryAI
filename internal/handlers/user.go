package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/service"
	"github.com/windoze95/saltybytes-api/internal/util"
)

// UserHandler is the handler for user-related requests.
type UserHandler struct {
	Service *service.UserService
}

// NewUserHandler is the constructor function for initializing a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{Service: userService}
}

// CreateUser creates a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser struct {
		Username  string `json:"username" binding:"required"`
		FirstName string `json:"first_name"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password fields are required"})
		return
	}

	// Validate username
	if err := h.Service.ValidateUsername(newUser.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email
	if err := h.Service.ValidateEmail(newUser.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password
	if err := h.Service.ValidatePassword(newUser.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user, err := h.Service.CreateUser(newUser.Username, newUser.FirstName, newUser.Email, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the user in
	tokenString, err := generateAuthToken(user.ID, h.Service.Cfg.Env.JwtSecretKey.Value())
	if err != nil {
		log.Printf("error: handlers.LoginUser: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenString, "message": "User signed up successfully", "user": user})
}

// LoginUser logs a user in.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var userCredentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	userResponse, err := h.Service.LoginUser(userCredentials.Username, userCredentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Log the user in
	tokenString, err := generateAuthToken(userResponse.ID, h.Service.Cfg.Env.JwtSecretKey.Value())
	if err != nil {
		log.Printf("error: handlers.LoginUser: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenString, "message": "User logged in successfully", "user": userResponse})
}

// generateAuthToken generates a JWT token for a user.
func generateAuthToken(userID uint, secretKey string) (string, error) {
	// Create a new token object, specifying signing method and the claims you would like it to contain.
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("generateAuthToken: %v ", err)
	}

	return tokenString, nil
}

// VerifyToken verifies a user's JWT token.
func (h *UserHandler) VerifyToken(c *gin.Context) {
	// Retrieve the user from the context
	user, _ := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"isAuthenticated": true, "user": user})
}

// GetUserByID fetches a user by ID.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Retrieve the user from the context
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserSettings fetches a user with settings.
func (h *UserHandler) GetUserSettings(c *gin.Context) {
	// Retrieve the user from the context
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": user.Settings})
}
