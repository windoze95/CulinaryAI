package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/service"
	"github.com/windoze95/culinaryai/internal/util"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{Service: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser struct {
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Recaptcha string `json:"recaptcha" binding:"required"`
	}

	// Read and log the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't read body"})
		return
	}

	// Log the request body
	log.Println(string(bodyBytes))

	// Restore the io.ReadCloser to its original state
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	log.Println("Validating reCAPTCHA")

	// Verify reCAPTCHA
	if err := h.Service.VerifyRecaptcha(newUser.Recaptcha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Validating username")

	// Validate username
	if err := h.Service.ValidateUsername(newUser.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Validating password")

	// Validate password
	if err := h.Service.ValidatePassword(newUser.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Hashing password")

	// Create user
	err := h.Service.CreateUser(newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("User created successfully")

	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var userCredentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.LoginUser(userCredentials.Username, userCredentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// // Create a new session
	// session := c.MustGet("session").(*sessions.Session)
	// session.Values["user_id"] = user.ID
	// session.Values["ip"] = c.ClientIP()
	// session.Values["user_agent"] = c.Request.UserAgent()

	// // Save the session
	// session.Save(c.Request, c.Writer)

	// c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(h.Service.Cfg.Env.JwtSecretKey.Value())
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not log in"})
		return
	}

	c.JSON(200, gin.H{"accessToken": tokenString, "message": "User logged in successfully", "user": user})
}

func (h *UserHandler) GetUserSettings(c *gin.Context) {
	// Retrieve the user from the context
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use the service to get and verify the settings
	isValid, err := h.Service.VerifyOpenAIKeyInUserSettings(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isValid": isValid, "user": user})
}

func (h *UserHandler) UpdateUserSettings(c *gin.Context) {
	// Retrieve the user from the context
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the new OpenAI key from the request body
	var newSettings struct {
		OpenAIKey string `json:"apikey"`
	}
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openAIKeyChanged, err := h.Service.UpdateUserSettings(user, newSettings.OpenAIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings: " + err.Error()})
		return
	}

	// This won't seem as redundant when more settings are added
	if openAIKeyChanged {
		c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "No changes made"})
	}
}
