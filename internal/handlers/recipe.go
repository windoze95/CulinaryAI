package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/repository"
	"github.com/windoze95/saltybytes-api/internal/service"
	"github.com/windoze95/saltybytes-api/internal/util"
)

type RecipeHandler struct {
	Service *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{Service: recipeService}
}

func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	recipeID := c.Param("recipe_id")

	recipe, err := h.Service.GetRecipeByID(recipeID)
	if err != nil {
		log.Printf("Error getting recipe: %v", err)
		switch e := err.(type) {
		case repository.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	// Retrieve the user from the context
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	log.Printf("User: %+v", user)
	log.Printf("User ID: %v", user.GuidingContent)

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

	go h.Service.CompleteRecipeGeneration(recipe, user)

	c.JSON(http.StatusOK, gin.H{"recipe": recipe, "message": "Generating recipe"})

	// go h.Service.CompleteRecipeGeneration(recipe, user)
}
