package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Recipe struct {
	gorm.Model
	Title                  string
	GeneratedRecipeVersion int `gorm:"type:int"`
	// GeneratedRecipe        openai.GeneratedRecipe `gorm:"-"`
	MainRecipeJSON     string `gorm:"type:text"`
	SubRecipesJSON     string `gorm:"type:text"`
	Tags               []Tag  `gorm:"many2many:recipe_tags;"`
	ImageURL           string
	GeneratedBy        *User `gorm:"foreignKey:GeneratedByUserID"`
	GeneratedByUserID  uint
	UserPrompt         string
	GuidingContentID   uint
	GuidingContentUID  uuid.UUID
	GuidingContent     *GuidingContent `gorm:"foreignKey:GuidingContentID"`
	GenerationComplete bool
}

type Tag struct {
	gorm.Model
	Hashtag string `gorm:"index:idx_hashtag;unique"`
}

// // SerializeGeneratedRecipe serializes the GeneratedRecipe field to a JSON string
// func (r *Recipe) SerializeGeneratedRecipe() error {
// 	// Set the current version
// 	r.GeneratedRecipeVersion = 1

// 	// Create an anonymous struct with only the fields you want to serialize
// 	tempStruct := struct {
// 		MainRecipe  MainRecipe   `json:"main_recipe"`
// 		SubRecipes  []MainRecipe `json:"sub_recipes"`
// 		DallEPrompt string       `json:"dall_e_prompt"`
// 	}{
// 		MainRecipe:  r.GeneratedRecipe.MainRecipe,
// 		SubRecipes:  r.GeneratedRecipe.SubRecipes,
// 		DallEPrompt: r.GeneratedRecipe.DallEPrompt,
// 	}

// 	generatedRecipeJSON, err := json.Marshal(tempStruct)
// 	if err != nil {
// 		return err
// 	}
// 	r.GeneratedRecipeJSON = string(generatedRecipeJSON)
// 	return nil
// }

// // DeserializeGeneratedRecipe deserializes the GeneratedRecipeJSON field back into the GeneratedRecipe struct
// func (r *Recipe) DeserializeGeneratedRecipe() error {
// 	// Use the version to determine how to deserialize GeneratedRecipe
// 	switch r.GeneratedRecipeVersion {
// 	case 1:
// 		// Deserialize directly into the GeneratedRecipe field, populating all its fields
// 		return json.Unmarshal([]byte(r.GeneratedRecipeJSON), &r.GeneratedRecipe)
// 	default:
// 		return fmt.Errorf("unsupported GeneratedRecipe version: %d", r.GeneratedRecipeVersion)
// 	}
// }
