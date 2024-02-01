package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// RecipeDef is a struct that represents the JSON schema that is passed to the OpenAI API for recipe generation using function calling.
type RecipeDef struct {
	Title             string         `json:"title" gorm:"column:title"`
	Ingredients       Ingredients    `json:"ingredients" gorm:"type:jsonb;column:ingredients"`
	Instructions      pq.StringArray `json:"instructions" gorm:"type:text[];column:instructions"`
	CookTime          int            `json:"cook_time" gorm:"column:cook_time"`
	ImagePrompt       string         `json:"image_prompt" gorm:"column:image_prompt"`
	Hashtags          []string       `json:"hashtags"` // Hashtags is shadowed by the Hashtags field in the Recipe model
	LinkedSuggestions pq.StringArray `json:"linked_recipe_suggestions" gorm:"type:text[];column:linked_recipe_suggestions"`
	// UnitSystem              UnitSystem   `json:"unit_system"`
}

// Scan is a GORM hook that scans jsonb into a RecipeDef.
func (j *RecipeDef) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := RecipeDef{}
	err := json.Unmarshal(bytes, &result)
	*j = RecipeDef(result)

	return err
}

// Value is a GORM hook that returns json value of a RecipeDef.
func (j RecipeDef) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Ingredient is a struct that represents an ingredient in a recipe.
type Ingredient struct {
	Name   string  `json:"name"`
	Unit   string  `json:"unit"`
	Amount float64 `json:"amount"`
}

// Ingredients is a slice of Ingredient.
// This is a workaround for GORM to embed a slice of structs into a JSONB field.
type Ingredients []Ingredient

// Scan is a GORM hook that scans jsonb into Ingredients.
func (j *Ingredients) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Ingredients{}
	err := json.Unmarshal(bytes, &result)
	*j = Ingredients(result)

	return err
}

// Value is a GORM hook that returns json value of Ingredients.
func (j Ingredients) Value() (driver.Value, error) {
	return json.Marshal(j)
}
