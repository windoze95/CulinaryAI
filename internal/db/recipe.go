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

func (db *RecipeDB) CreateRecipe(recipe *models.Recipe) error {
	return db.DB.Create(recipe).Error
}

// func (db *RecipeDB) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return db.DB.Model(&models.Recipe{}).Where("id = ?", id).Update(field, value).Error
// }

func (db *RecipeDB) UpdateRecipeImageURL(recipe *models.Recipe, imageURL string) error {
	return db.DB.Model(recipe).Update("ImageURL", imageURL).Error
}

func (db *RecipeDB) UpdateRecipeGenerationStatus(recipe *models.Recipe, isComplete bool) error {
	return db.DB.Model(recipe).Update("GenerationComplete", isComplete).Error
}

func (db *RecipeDB) UpdateFullRecipeJSON(recipe *models.Recipe) error {
	return db.DB.Model(recipe).Update("FullRecipeJSON", recipe.FullRecipeJSON).Error
}

func (db *RecipeDB) FindTagByName(tagName string) (*models.Tag, error) {
	var tag models.Tag
	err := db.DB.Where("name = ?", tagName).First(&tag).Error
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
