package repository

import (
	"github.com/windoze95/culinaryai/internal/db"
	"github.com/windoze95/culinaryai/internal/models"
)

type RecipeRepository struct {
	recipeDB *db.RecipeDB
}

func NewRecipeRepository(recipeDB *db.RecipeDB) *RecipeRepository {
	return &RecipeRepository{recipeDB: recipeDB}
}

func (r *RecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	return r.recipeDB.CreateRecipe(recipe)
}

// func (r *RecipeRepository) UpdateRecipeFieldByID(id uint, field string, value interface{}) error {
// 	return r.recipeDB.UpdateRecipeFieldByID(id, field, value)
// }

func (r *RecipeRepository) UpdateRecipeImageURL(recipe *models.Recipe, imageURL string) error {
	return r.recipeDB.UpdateRecipeImageURL(recipe, imageURL)
}

func (r *RecipeRepository) UpdateRecipeGenerationStatus(recipe *models.Recipe, isComplete bool) error {
	return r.recipeDB.UpdateRecipeGenerationStatus(recipe, isComplete)
}

func (r *RecipeRepository) FindTagByName(tagName string) (*models.Tag, error) {
	return r.recipeDB.FindTagByName(tagName)
}

func (r *RecipeRepository) CreateTag(tag *models.Tag) error {
	return r.recipeDB.CreateTag(tag)
}

func (r *RecipeRepository) UpdateRecipeTagsAssociation(recipe *models.Recipe, tags []models.Tag) error {
	return r.recipeDB.UpdateRecipeTagsAssociation(recipe, tags)
}
