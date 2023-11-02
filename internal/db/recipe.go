package db

import (
	"log"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/windoze95/saltybytes-api/internal/models"
)

type RecipeDB struct {
	DB *gorm.DB
}

func NewRecipeDB(gormDB *gorm.DB) *RecipeDB {
	return &RecipeDB{DB: gormDB}
}

// func (db *RecipeDB) GetRecipeByID(id string) (*models.Recipe, error) {
// 	log.Printf("About to query database for recipe with ID: %s", id)
// 	var recipe models.Recipe
// 	err := db.DB.Preload("GuidingContent").
// 		Preload("Tags").
// 		Preload("GeneratedBy", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("Username")
// 		}).
// 		Where("id = ?", id).
// 		First(&recipe).Error
// 	log.Printf("Query complete. Recipe retrieved: %+v, Error: %v", recipe, err)
// 	return &recipe, err
// }

// func (db *RecipeDB) GetRecipeByID(id string) (*models.Recipe, error) {
// 	var recipe models.Recipe
// 	err := db.DB.Preload("GuidingContent").
// 		Preload("Tags").
// 		// Preload("GeneratedBy").
// 		Where("id = ?", id).
// 		First(&recipe).Error
// 	log.Printf("Query complete. Recipe retrieved: %+v, Error: %v", recipe, err)
// 	return &recipe, err
// }

func (db *RecipeDB) GetRecipeByID(recipeID string) (*models.Recipe, error) {
	var recipe models.Recipe
	err := db.DB.Preload("GuidingContent").
		Preload("Tags").
		Preload("GeneratedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("Username") // Select only Username
		}).
		Where("id = ?", recipeID).
		First(&recipe).Error
	if err != nil {
		log.Printf("Error retrieving recipe: %v", err)
		return nil, err
	}
	// log.Printf("Query complete. username retrieved: %+v, Error: %v", recipe.GeneratedBy, err)
	return &recipe, err
}

func (db *RecipeDB) GetChatHistoryByRecipeID(recipeID uint) (*models.RecipeChatHistory, error) {
	var chatHistory models.RecipeChatHistory
	err := db.DB.Where("RecipeID = ?", recipeID).
		First(&chatHistory).Error
	if err != nil {
		log.Printf("Error retrieving chat history: %v", err)
		return nil, err
	}
	return &chatHistory, nil
}

func (db *RecipeDB) CreateRecipe(recipe *models.Recipe) error {
	// Start a new transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Omit("GeneratedBy").
		Omit("GuidingContent").
		Create(recipe).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating recipe: %v", err)
		return err
	}

	return tx.Commit().Error
}

// Creating chat history might be restrained to a recipe for strict association requirements like setting recipeID

// func (db *RecipeDB) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return db.DB.Model(&models.Recipe{}).Where("id = ?", id).Update(field, value).Error
// }

func (db *RecipeDB) UpdateRecipeTitle(recipe *models.Recipe, title string) error {
	err := db.DB.Model(recipe).
		Update("Title", title).Error
	if err != nil {
		log.Printf("Error updating recipe title: %v", err)
		return err
	}
	return nil
}

func (db *RecipeDB) UpdateRecipeImageURL(recipe *models.Recipe, imageURL string) error {
	err := db.DB.Model(recipe).
		Update("ImageURL", imageURL).Error
	if err != nil {
		log.Printf("Error updating recipe image URL: %v", err)
		return err
	}
	return nil
}

func (db *RecipeDB) UpdateRecipeGenerationStatus(recipe *models.Recipe, isComplete bool) error {
	err := db.DB.Model(recipe).
		Update("GenerationComplete", isComplete).Error
	if err != nil {
		log.Printf("Error updating recipe generation status: %v", err)
		return err
	}
	return nil
}

// func (db *RecipeDB) UpdateGeneratedRecipeJSON(recipe *models.Recipe) error {
// 	return db.DB.Model(recipe).Update("GeneratedRecipeJSON", recipe.GeneratedRecipeJSON).Error
// }

//	func (db *RecipeDB) UpdateRecipeCoreFields(recipe *models.Recipe) error {
//		err := db.DB.Model(recipe).
//			Updates(map[string]interface{}{
//				"Title":          recipe.Title,
//				"MainRecipeJSON": recipe.MainRecipeJSON,
//				"SubRecipesJSON": recipe.SubRecipesJSON,
//				"ImagePrompt":        recipe.ImagePrompt,
//				// "GeneratedRecipeVersion": recipe.GeneratedRecipeVersion,
//			}).Error
//		if err != nil {
//			log.Printf("Error updating recipe core fields: %v", err)
//			return err
//		}
//		return nil
//	}
func (db *RecipeDB) UpdateRecipeCoreFields(recipe *models.Recipe) error {
	// Start a new transaction.
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update core fields of the recipe.
	err := tx.Model(recipe).
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

	// Update the chat history.
	if recipe.ChatHistory != nil {
		err = tx.Model(recipe.ChatHistory).
			Updates(map[string]interface{}{
				"MessagesJSON": recipe.ChatHistory.MessagesJSON,
			}).Error
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating recipe chat history: %v", err)
			return err
		}
	}

	// Commit the transaction if all updates succeed.
	err = tx.Commit().Error
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (db *RecipeDB) FindTagByName(tagName string) (*models.Tag, error) {
	var tag models.Tag
	err := db.DB.Where("Hashtag = ?", tagName).
		First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (db *RecipeDB) CreateTag(tag *models.Tag) error {
	err := db.DB.Create(tag).Error
	if err != nil {
		log.Printf("Error creating tag: %v", err)
		return err
	}
	return nil
}

func (db *RecipeDB) UpdateRecipeTagsAssociation(recipe *models.Recipe, tags []models.Tag) error {
	err := db.DB.Model(&recipe).
		Association("Hashtags").
		Replace(tags).Error
	if err != nil {
		log.Printf("Error updating recipe tags association: %v", err)
		return err
	}
	return nil
}
