package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/service"
)

type RecipeHandler struct {
	Service *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{Service: recipeService}
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	// Retrieve the user from the context
	user, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the request body for the user's prompt
	var request struct {
		UserPrompt string `json:"userPrompt"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	recipe, err := h.Service.CreateRecipe(user, request.UserPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe, "message": "Initial recipe saved, generating full recipe"})

	go h.Service.CompleteRecipeGeneration(recipe, user)
}
