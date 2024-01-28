package repository

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/windoze95/saltybytes-api/internal/models"
)

// RecipeRepository is a repository for interacting with recipes.
type RecipeRepository struct {
	DB *gorm.DB
}

// NewRecipeRepository creates a new RecipeRepository.
func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{DB: db}
}

// GetRecipeByID retrieves a recipe by its ID.
func (r *RecipeRepository) GetRecipeByID(recipeID uint) (*models.Recipe, error) {
	var recipe models.Recipe

	err := r.DB.Preload("Hashtags").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("Username") // Select only Username
		}).
		Where("id = ?", recipeID).
		First(&recipe).Error
	if err != nil {
		log.Printf("Error retrieving recipe: %v", err)

		if gorm.IsRecordNotFoundError(err) {
			return nil, NotFoundError{message: "Recipe not found"}
		}

		return nil, err
	}

	return &recipe, nil
}

// GetHistoryByID retrieves a recipe history by its ID.
func (r *RecipeRepository) GetHistoryByID(historyID uint) (*models.RecipeHistory, error) {
	var history *models.RecipeHistory

	err := r.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(history, historyID).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

// GetMessagesAfterID retrieves messages belonging to a specific RecipeHistory
// and having an ID greater than a given value.
func (r *RecipeRepository) GetMessagesAfterID(historyID uint, afterID uint) ([]models.RecipeHistoryMessage, error) {
	var messages []models.RecipeHistoryMessage

	result := r.DB.Where("recipe_chat_history_id = ? AND id > ?", historyID, afterID).
		Order("id ASC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

// CreateRecipe creates a new recipe.
func (r *RecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	// Start a new transaction
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Create(recipe).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating recipe: %v", err)
		return err
	}

	return tx.Commit().Error
}

// DeleteRecipe deletes a recipe.
func (r *RecipeRepository) DeleteRecipe(recipeID uint) error {
	err := r.DB.Delete(&models.Recipe{}, recipeID).Error
	if err != nil {
		log.Printf("Error deleting recipe: %v", err)
	}
	return err
}

// UpdateRecipeTitle updates the title of a recipe.
func (r *RecipeRepository) UpdateRecipeTitle(recipe *models.Recipe, title string) error {
	err := r.DB.Model(recipe).
		Update("Title", title).Error
	if err != nil {
		log.Printf("Error updating recipe title: %v", err)
	}
	return err
}

// UpdateRecipeImageURL updates the image URL of a recipe.
func (r *RecipeRepository) UpdateRecipeImageURL(recipeID uint, imageURL string) error {
	err := r.DB.Model(&models.Recipe{}).
		Where("id = ?", recipeID).
		Update("ImageURL", imageURL).Error
	if err != nil {
		log.Printf("Error updating recipe image URL: %v", err)
	}
	return err
}

// UpdateRecipeDef updates the core fields of a recipe and appends new messages to the history.
//
// Core fields: "Title", "Ingredients", "Instructions", "CookTime", "LinkedRecipeSuggestions", "ImagePrompt"
func (r *RecipeRepository) UpdateRecipeDef(recipe *models.Recipe, newRecipeHistoryMessage models.RecipeHistoryMessage) error {
	// Start a new transaction.
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update core fields of the recipe.
	err := tx.Model(&models.Recipe{}).
		Where("id = ?", recipe.ID).
		Updates(map[string]interface{}{
			"Title":           recipe.Title,
			"Ingredients":     recipe.Ingredients,
			"Instructions":    recipe.Instructions,
			"CookTime":        recipe.CookTime,
			"LinkSuggestions": recipe.LinkSuggestions,
			"ImagePrompt":     recipe.ImagePrompt,
		}).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating recipe core fields: %v", err)
		return err
	}

	// Check if HistoryID is set in the Recipe
	if recipe.HistoryID == 0 {
		tx.Rollback()
		err = errors.New("recipe history ID not set in recipe")
		log.Printf("Error: %v", err)
		return err
	}

	newRecipeHistoryMessage.RecipeHistoryID = recipe.HistoryID

	// Insert the new message into the database
	err = tx.Create(&newRecipeHistoryMessage).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating new recipe history message: %v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Printf("Error committing transaction in UpdateRecipeCoreFields: %v", err)
		return err
	}

	return nil
}

// FindTagByName finds a tag by its name.
func (r *RecipeRepository) FindTagByName(tagName string) (*models.Tag, error) {
	var tag models.Tag
	err := r.DB.Where("Hashtag = ?", tagName).
		First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// CreateTag creates a new tag.
func (r *RecipeRepository) CreateTag(tag *models.Tag) error {
	err := r.DB.Create(tag).Error
	if err != nil {
		log.Printf("Error creating tag: %v", err)
	}
	return err
}

// UpdateRecipeTagsAssociation updates the tags associated with a recipe.
func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipeID uint, newTags []models.Tag) error {
	var recipe models.Recipe
	result := r.DB.First(&recipe, recipeID)
	if result.Error != nil {
		return result.Error
	}

	// Replace existing associations with new tags
	err := r.DB.Model(&recipe).
		Association("Hashtags").
		Replace(newTags).Error
	if err != nil {
		return err
	}

	return nil
}
