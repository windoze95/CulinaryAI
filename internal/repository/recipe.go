package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/windoze95/saltybytes-api/internal/models"
)

// type RecipeRepository struct {
// 	RecipeDB *db.RecipeDB
// }

type RecipeRepository struct {
	DB *gorm.DB
}

// func NewRecipeRepository(recipeDB *db.RecipeDB) *RecipeRepository {
// 	return &RecipeRepository{RecipeDB: recipeDB}
// }

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{DB: db}
}

func (r *RecipeRepository) GetRecipeByID(recipeID uint) (*models.Recipe, error) {
	// recipe, err := r.RecipeDB.GetRecipeByID(recipeID)
	var recipe models.Recipe
	err := r.DB.Preload("GuidingContent").
		Preload("Tags").
		Preload("GeneratedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("Username") // Select only Username
		}).
		Where("id = ?", recipeID).
		First(&recipe).Error
	if err != nil {
		log.Printf("Error retrieving recipe: %v", err)
		// 	return nil, err
		// }
		// log.Printf("Query complete. username retrieved: %+v, Error: %v", recipe.GeneratedBy, err)
		// if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, NotFoundError{message: "Recipe not found"}
		}
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) GetChatHistoryByRecipeID(recipeID uint) (*models.RecipeChatHistory, error) {
	var chatHistory models.RecipeChatHistory
	err := r.DB.Where("RecipeID = ?", recipeID).
		First(&chatHistory).Error
	if err != nil {
		log.Printf("Error retrieving chat history: %v", err)
		return nil, err
	}
	return &chatHistory, nil
}

func (r *RecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	// return r.RecipeDB.CreateRecipe(recipe)
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

// func (r *RecipeRepository) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return r.RecipeDB.UpdateRecipeFieldByID(id, field, value)
// }

func (r *RecipeRepository) UpdateRecipeTitle(recipe *models.Recipe, title string) error {
	// return r.RecipeDB.UpdateRecipeTitle(recipe, title)
	err := r.DB.Model(recipe).
		Update("Title", title).Error
	if err != nil {
		log.Printf("Error updating recipe title: %v", err)
	}
	return err
}

func (r *RecipeRepository) UpdateRecipeImageURL(recipeID uint, imageURL string) error {
	// return r.RecipeDB.UpdateRecipeImageURL(recipe, imageURL)
	err := r.DB.Model(&models.Recipe{}).
		Where("id = ?", recipeID).
		Update("ImageURL", imageURL).Error
	if err != nil {
		log.Printf("Error updating recipe image URL: %v", err)
	}
	return err
}

func (r *RecipeRepository) UpdateRecipeGenerationStatus(recipeID uint, isComplete bool) error {
	// return r.RecipeDB.UpdateRecipeGenerationStatus(recipe, isComplete)
	err := r.DB.Model(&models.Recipe{}).
		Where("id = ?", recipeID).
		Update("GenerationComplete", isComplete).Error
	if err != nil {
		log.Printf("Error updating recipe generation status: %v", err)
	}
	return err
}

func (r *RecipeRepository) UpdateRecipeCoreFields(recipe *models.Recipe, newRecipeChatHistoryMessages []string) error {
	// return r.RecipeDB.UpdateRecipeCoreFields(recipe, newRecipeChatHistoryMessages)
	// Start a new transaction.
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update core fields of the recipe.
	err := tx.Model(&models.Recipe{}).
		Where("id = ?", recipe.ID).
		Updates(map[string]interface{}{
			"Title":          recipe.Title,
			"MainRecipeJSON": recipe.MainRecipeJSON,
			"SubRecipesJSON": recipe.SubRecipesJSON,
			"ImagePrompt":    recipe.ImagePrompt,
		}).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating recipe core fields: %v", err)
		return err
	}

	// Append new messages to the chat history.
	// for _, message := range newRecipeChatHistoryMessages {
	// 	// Properly escape the JSON string for SQL query
	// 	escapedMessage := strings.ReplaceAll(message, "'", "''")

	// 	err = tx.Exec(`UPDATE recipe_chat_histories SET messages_json = array_append(messages_json, ?) WHERE id = ?`, escapedMessage, recipe.ChatHistory.ID).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 		log.Printf("Error appending message to recipe chat history: %v", err)
	// 		return err
	// 	}
	// }
	if len(newRecipeChatHistoryMessages) > 0 {
		// Convert the new messages into a PostgreSQL array literal
		// newMessagesPGArray := "{" + strings.Join(newRecipeChatHistoryMessages, ",") + "}"
		err = tx.Exec("UPDATE recipe_chat_histories SET messages_json = array_cat(messages_json, ?) WHERE id = ?", newRecipeChatHistoryMessages, recipe.ChatHistory.ID).Error
		// err = tx.Exec(`UPDATE recipe_chat_histories SET messages_json = array_cat(messages_json, ?) WHERE id = ?`, newMessagesPGArray, recipe.ChatHistory.ID).Error
		if err != nil {
			tx.Rollback()
			log.Printf("Error appending messages to recipe chat history: %v", err)
			return err
		}
	}

	// // Commit the transaction if all updates succeed.
	// err = tx.Commit().Error
	// if err != nil {
	// 	log.Printf("Error committing transaction: %v", err)
	// 	return err
	// }

	return tx.Commit().Error
}

func (r *RecipeRepository) FindTagByName(tagName string) (*models.Tag, error) {
	// return r.RecipeDB.FindTagByName(tagName)
	var tag models.Tag
	err := r.DB.Where("Hashtag = ?", tagName).
		First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *RecipeRepository) CreateTag(tag *models.Tag) error {
	// return r.RecipeDB.CreateTag(tag)
	err := r.DB.Create(tag).Error
	if err != nil {
		log.Printf("Error creating tag: %v", err)
	}
	return err
}

func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipeID uint, tags []models.Tag) error {
	// return r.RecipeDB.UpdateRecipeTagsAssociation(recipe, tags)
	err := r.DB.Model(&models.Recipe{}).
		Where("id = ?", recipeID).
		Association("Hashtags").
		Replace(tags).Error
	if err != nil {
		log.Printf("Error updating recipe tags association: %v", err)
	}
	return err
}
