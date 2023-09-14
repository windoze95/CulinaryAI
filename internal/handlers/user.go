package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/windoze95/culinaryai/internal/service"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{Service: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	err := h.Service.CreateUser(newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	user, err := h.Service.Login(userCredentials.Username, userCredentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*sessions.Session)
	session.Values["user_id"] = user.ID
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

func (h *UserHandler) GetSettings(c *gin.Context) {
	// Retrieve the user from the context
	user, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use the service to get and verify the settings
	isValid, err := h.Service.VerifyOpenAIKeyInSettings(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the settings modal template
	c.HTML(http.StatusOK, "settings.tmpl", gin.H{"isValid": isValid, "user": user})
}

func (h *UserHandler) UpdateUserSettings(c *gin.Context) {
	// Retrieve the user from the context
	user, err := getUserFromContext(c)
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
