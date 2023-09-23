package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Recipe struct {
	gorm.Model
	Title              string
	FullRecipeVersion  int        `gorm:"type:int"`
	FullRecipe         FullRecipe `gorm:"-"`
	FullRecipeJSON     string     `gorm:"type:text"`
	Tags               []Tag      `gorm:"many2many:recipe_tags;"`
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

type Ingredient struct {
	Name   string  `json:"name"`
	Unit   string  `json:"unit"`
	Amount float64 `json:"amount"`
}

type MainRecipe struct {
	RecipeName   string       `json:"recipe_name"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
	TimeToCook   int          `json:"time_to_cook"`
}

type FullRecipe struct {
	Title       string       `json:"title"`
	MainRecipe  MainRecipe   `json:"main_recipe"`
	SubRecipes  []MainRecipe `json:"sub_recipes"`
	DallEPrompt string       `json:"dall_e_prompt"`
	UnitSystem  string       `json:"unit_system"` // This field will not be serialized, but will be deserialized
	Hashtags    []string     `json:"hashtags"`    // This field will not be serialized, but will be deserialized
	// UnitSystem  string       `json:"-"`
	// Hashtags    []string     `json:"-"`
	// UnitSystem  string       `json:"unit_system"`
	// Hashtags    []string     `json:"hashtags"`
}

// SerializeFullRecipe serializes the FullRecipe field to a JSON string
func (r *Recipe) SerializeFullRecipe() error {
	// Set the current version
	r.FullRecipeVersion = 1

	// Create an anonymous struct with only the fields you want to serialize
	tempStruct := struct {
		MainRecipe  MainRecipe   `json:"main_recipe"`
		SubRecipes  []MainRecipe `json:"sub_recipes"`
		DallEPrompt string       `json:"dall_e_prompt"`
	}{
		MainRecipe:  r.FullRecipe.MainRecipe,
		SubRecipes:  r.FullRecipe.SubRecipes,
		DallEPrompt: r.FullRecipe.DallEPrompt,
	}

	fullRecipeJSON, err := json.Marshal(tempStruct)
	if err != nil {
		return err
	}
	r.FullRecipeJSON = string(fullRecipeJSON)
	return nil
}

// DeserializeFullRecipe deserializes the FullRecipeJSON field back into the FullRecipe struct
func (r *Recipe) DeserializeFullRecipe() error {
	// Use the version to determine how to deserialize FullRecipe
	switch r.FullRecipeVersion {
	case 1:
		// Deserialize directly into the FullRecipe field, populating all its fields
		return json.Unmarshal([]byte(r.FullRecipeJSON), &r.FullRecipe)
	default:
		return fmt.Errorf("unsupported FullRecipe version: %d", r.FullRecipeVersion)
	}
}
