package repository

import (
	"errors"
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
		Preload("Hashtags").
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

//	func (r *RecipeRepository) GetChatHistoryByRecipeID(chatHistoryID uint) (*models.RecipeChatHistory, error) {
//		var chatHistory models.RecipeChatHistory
//		err := r.DB.Where("id = ?", chatHistoryID).
//			First(&chatHistory).Error
//		if err != nil {
//			log.Printf("Error retrieving chat history: %v", err)
//			return nil, err
//		}
//		return &chatHistory, nil
//	}
func (r *RecipeRepository) GetChatHistoryByID(chatHistoryID uint) (*models.RecipeChatHistory, error) {
	var chatHistory models.RecipeChatHistory
	// err := r.DB.Order("created_at ASC").
	// 	Find(&chatHistory, "recipe_chat_history_id = ?", chatHistoryID).Error
	err := r.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&chatHistory, chatHistoryID).Error
	if err != nil {
		return nil, err
	}
	// result := r.DB.First(&chatHistory, chatHistoryID)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	log.Printf("Repo: Chat history fetched: %+v", chatHistory)
	log.Printf("Repo: Chat history messages main recipe name: %+v", chatHistory.Messages[0].GeneratedResponse.MainRecipe.RecipeName)

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

func (r *RecipeRepository) UpdateRecipeCoreFields(recipe *models.Recipe, newRecipeChatHistoryMessage models.RecipeChatHistoryMessage) error {
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

	// if len(newRecipeChatHistoryMessages) > 0 {
	// 	// // Convert newRecipeChatHistoryMessages to a PostgreSQL-compatible array format
	// 	// pgArray := "{" + strings.Join(newRecipeChatHistoryMessages, ",") + "}"

	// 	// Use a parameterized query to safely append messages
	// 	err = tx.Exec("UPDATE recipe_chat_histories SET messages_json = array_cat(messages_json, ?::text[]) WHERE id = ?", pq.Array(newRecipeChatHistoryMessages), recipe.ChatHistory.ID).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 		log.Printf("Error appending messages to recipe chat history: %v", err)
	// 		return err
	// 	}
	// }

	// Check if ChatHistoryID is set in the Recipe
	if recipe.ChatHistoryID == 0 {
		tx.Rollback()
		err = errors.New("chat history ID not set in recipe")
		log.Printf("Error: %v", err)
		return err
	}

	newRecipeChatHistoryMessage.RecipeChatHistoryID = recipe.ChatHistoryID

	// Insert the new message into the database
	err = tx.Create(&newRecipeChatHistoryMessage).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating new recipe chat history message: %v", err)
		return err
	}

	// // Iterate over the slice of new messages
	// for i := range newRecipeChatHistoryMessages {
	// 	// Set the foreign key to link each new message to the specific RecipeChatHistory
	// 	newRecipeChatHistoryMessages[i].RecipeChatHistoryID = recipe.ChatHistoryID

	// 	// Directly insert each new message into the database
	// 	err = tx.Create(&newRecipeChatHistoryMessages[i]).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if len(newRecipeChatHistoryMessages) > 0 {
	// 	// Convert the new messages into a PostgreSQL array literal
	// 	// newMessagesPGArray := "{" + strings.Join(newRecipeChatHistoryMessages, ",") + "}"
	// 	err = tx.Exec("UPDATE recipe_chat_histories SET messages_json = array_cat(messages_json, ?) WHERE id = ?", newRecipeChatHistoryMessages, recipe.ChatHistory.ID).Error
	// 	// err = tx.Exec(`UPDATE recipe_chat_histories SET messages_json = array_cat(messages_json, ?) WHERE id = ?`, newMessagesPGArray, recipe.ChatHistory.ID).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 		log.Printf("Error appending messages to recipe chat history: %v", err)
	// 		return err
	// 	}
	// }

	// // Commit the transaction if all updates succeed.
	// err = tx.Commit().Error
	// if err != nil {
	// 	log.Printf("Error committing transaction: %v", err)
	// 	return err
	// }

	err = tx.Commit().Error
	if err != nil {
		log.Printf("Error committing transaction in UpdateRecipeCoreFields: %v", err)
		return err
	}

	return nil
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

// func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipeID uint, tags []models.Tag) error {
// 	// return r.RecipeDB.UpdateRecipeTagsAssociation(recipe, tags)
// 	// Debug logs
// 	log.Printf("Updating tags for recipe ID: %d", recipeID)
// 	for _, tag := range tags {
// 		log.Printf("Tag ID: %d, Hashtag: %s", tag.ID, tag.Hashtag)
// 	}
// 	err := r.DB.Model(&models.Recipe{}).
// 		Where("id = ?", recipeID).
// 		Association("Hashtags").
// 		Replace(tags).Error
// 	if err != nil {
// 		log.Printf("Error updating recipe tags association: %v", err)
// 	}
// 	return err
// }

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

// func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipeID uint, tags []models.Tag) error {
// 	err := r.DB.Debug().Model(&models.Recipe{}).
// 		Where("id = ?", recipeID).
// 		Association("Hashtags").
// 		Replace(tags).Error
// 	if err != nil {
// 		log.Printf("Error updating recipe tags association: %v", err)
// 	}
// 	return err
// }
