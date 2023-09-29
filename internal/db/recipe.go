package db

import (
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/windoze95/culinaryai/internal/models"
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

func (db *RecipeDB) GetRecipeByID(id string) (*models.Recipe, error) {
	var recipe models.Recipe
	err := db.DB.Preload("GuidingContent").
		Preload("Tags").
		Preload("GeneratedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username") // Select only ID and Username
		}).
		Where("id = ?", id).
		First(&recipe).Error
	// log.Printf("Query complete. username retrieved: %+v, Error: %v", recipe.GeneratedBy, err)
	return &recipe, err
}

func (db *RecipeDB) CreateRecipe(recipe *models.Recipe) error {
	return db.DB.Omit("GeneratedBy").Omit("GuidingContent").Create(recipe).Error
}

// func (db *RecipeDB) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return db.DB.Model(&models.Recipe{}).Where("id = ?", id).Update(field, value).Error
// }

func (db *RecipeDB) UpdateRecipeTitle(recipe *models.Recipe, title string) error {
	return db.DB.Model(recipe).Update("Title", title).Error
}

func (db *RecipeDB) UpdateRecipeImageURL(recipe *models.Recipe, imageURL string) error {
	return db.DB.Model(recipe).Update("ImageURL", imageURL).Error
}

func (db *RecipeDB) UpdateRecipeGenerationStatus(recipe *models.Recipe, isComplete bool) error {
	return db.DB.Model(recipe).Update("GenerationComplete", isComplete).Error
}

// func (db *RecipeDB) UpdateFullRecipeJSON(recipe *models.Recipe) error {
// 	return db.DB.Model(recipe).Update("FullRecipeJSON", recipe.FullRecipeJSON).Error
// }

func (db *RecipeDB) UpdateFullRecipeJSON(recipe *models.Recipe) error {
	return db.DB.Model(recipe).Updates(map[string]interface{}{
		"FullRecipeJSON":    recipe.FullRecipeJSON,
		"FullRecipeVersion": recipe.FullRecipeVersion,
	}).Error
}

func (db *RecipeDB) FindTagByName(tagName string) (*models.Tag, error) {
	var tag models.Tag
	err := db.DB.Where("Hashtag = ?", tagName).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (db *RecipeDB) CreateTag(tag *models.Tag) error {
	return db.DB.Create(tag).Error
}

func (db *RecipeDB) UpdateRecipeTagsAssociation(recipe *models.Recipe, tags []models.Tag) error {
	return db.DB.Model(&recipe).Association("Tags").Replace(tags).Error
}
