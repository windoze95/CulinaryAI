package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// RecipeDef is a struct that represents the JSON schema that is passed to the OpenAI API for recipe generation using function calling.
type RecipeDef struct {
	Title        string        `json:"title"`
	Ingredients  []*Ingredient `json:"ingredients"`
	Instructions []string      `json:"instructions"`
	CookTime     int           `json:"cook_time"`
	ImagePrompt  string        `json:"image_prompt"`
	// UnitSystem              UnitSystem   `json:"unit_system"`
	Hashtags                []string `json:"hashtags"`
	LinkedRecipeSuggestions []string `json:"linked_recipe_suggestions"`
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

// Scan is a GORM hook that scans jsonb into a Ingredient.
func (j *Ingredient) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Ingredient{}
	err := json.Unmarshal(bytes, &result)
	*j = Ingredient(result)

	return err
}

// Value is a GORM hook that returns json value of a Ingredient.
func (j Ingredient) Value() (driver.Value, error) {
	return json.Marshal(j)
}
