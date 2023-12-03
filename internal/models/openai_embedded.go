package models

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
