package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/windoze95/culinaryai/internal/db"
	"github.com/windoze95/culinaryai/internal/models"
)

type RecipeRepository struct {
	RecipeDB *db.RecipeDB
}

func NewRecipeRepository(recipeDB *db.RecipeDB) *RecipeRepository {
	return &RecipeRepository{RecipeDB: recipeDB}
}

func (r *RecipeRepository) GetRecipeByID(recipeID string) (*models.Recipe, error) {
	recipe, err := r.RecipeDB.GetRecipeByID(recipeID)
	if err != nil {
		log.Printf("Repository error: %v", err)
		if gorm.IsRecordNotFoundError(err) {
			return nil, NotFoundError{message: "Recipe not found"}
		}
		return nil, err
	}
	return recipe, nil
}

func (r *RecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	return r.RecipeDB.CreateRecipe(recipe)
}

// func (r *RecipeRepository) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return r.RecipeDB.UpdateRecipeFieldByID(id, field, value)
// }

func (r *RecipeRepository) UpdateRecipeTitle(recipe *models.Recipe, title string) error {
	return r.RecipeDB.UpdateRecipeTitle(recipe, title)
}

func (r *RecipeRepository) UpdateRecipeImageURL(recipe *models.Recipe, imageURL string) error {
	return r.RecipeDB.UpdateRecipeImageURL(recipe, imageURL)
}

func (r *RecipeRepository) UpdateRecipeGenerationStatus(recipe *models.Recipe, isComplete bool) error {
	return r.RecipeDB.UpdateRecipeGenerationStatus(recipe, isComplete)
}

func (r *RecipeRepository) UpdateFullRecipeJSON(recipe *models.Recipe) error {
	return r.RecipeDB.UpdateFullRecipeJSON(recipe)
}

func (r *RecipeRepository) FindTagByName(tagName string) (*models.Tag, error) {
	return r.RecipeDB.FindTagByName(tagName)
}

func (r *RecipeRepository) CreateTag(tag *models.Tag) error {
	return r.RecipeDB.CreateTag(tag)
}

func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipe *models.Recipe, tags []models.Tag) error {
	return r.RecipeDB.UpdateRecipeTagsAssociation(recipe, tags)
}
