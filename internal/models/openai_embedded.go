package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type FunctionCallArgument struct {
	Title       string            `json:"title"`
	MainRecipe  GeneratedRecipe   `json:"main_recipe"`
	SubRecipes  []GeneratedRecipe `json:"sub_recipes"`
	ImagePrompt string            `json:"image_prompt"`
	UnitSystem  string            `json:"unit_system"`
	Hashtags    []string          `json:"hashtags"`
	// ChatContext string   `json:"chat_context"`
	// UnitSystem  string       `json:"-"`
	// Hashtags    []string     `json:"-"`
	// UnitSystem  string       `json:"unit_system"`
	// Hashtags    []string     `json:"hashtags"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *FunctionCallArgument) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := FunctionCallArgument{}
	err := json.Unmarshal(bytes, &result)
	*j = FunctionCallArgument(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j FunctionCallArgument) Value() (driver.Value, error) {
	// if len(j) == 0 {
	// 	return nil, nil
	// }
	// return json.RawMessage(j).MarshalJSON()
	return json.Marshal(j)
}

type Ingredient struct {
	Name   string  `json:"name"`
	Unit   string  `json:"unit"`
	Amount float64 `json:"amount"`
}

// type Recipe struct {
type GeneratedRecipe struct {
	RecipeName   string       `json:"recipe_name"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
	TimeToCook   float64      `json:"time_to_cook"`
}
