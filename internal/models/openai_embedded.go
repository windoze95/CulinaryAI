package models

type FunctionCallArgument struct {
	Title       string            `json:"title"`
	MainRecipe  GeneratedRecipe   `gorm:"type:json" json:"main_recipe"`
	SubRecipes  []GeneratedRecipe `gorm:"type:json" json:"sub_recipes"`
	ImagePrompt string            `json:"image_prompt"`
	UnitSystem  string            `json:"unit_system"` // This field will not be serialized, but will be deserialized
	Hashtags    []string          `json:"hashtags"`    // This field will not be serialized, but will be deserialized
	// ChatContext string   `json:"chat_context"`
	// UnitSystem  string       `json:"-"`
	// Hashtags    []string     `json:"-"`
	// UnitSystem  string       `json:"unit_system"`
	// Hashtags    []string     `json:"hashtags"`
}

type Ingredient struct {
	Name   string  `json:"name"`
	Unit   string  `json:"unit"`
	Amount float64 `json:"amount"`
}

// type Recipe struct {
type GeneratedRecipe struct {
	RecipeName   string       `json:"recipe_name"`
	Ingredients  []Ingredient `gorm:"type:json" json:"ingredients"`
	Instructions []string     `json:"instructions"`
	TimeToCook   float64      `json:"time_to_cook"`
}
