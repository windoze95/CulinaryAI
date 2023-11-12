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
	Hashtags           []Tag  `gorm:"many2many:recipe_tags;"`
	ImagePrompt        string
	ImageURL           string
	GeneratedBy        *User `gorm:"foreignKey:GeneratedByUserID"`
	GeneratedByUserID  uint
	InitialPrompt      string
	GuidingContentID   uint
	GuidingContentUID  uuid.UUID
	GuidingContent     *GuidingContent    `gorm:"foreignKey:GuidingContentID"`
	ChatHistoryID      uint               `gorm:"unique;index"`
	ChatHistory        *RecipeChatHistory `gorm:"foreignKey:ChatHistoryID"`
	SpinOnRecipeID     *uint
	SpinOnRecipe       *Recipe `gorm:"foreignKey:SpinOnRecipeID"`
	GenerationComplete bool
}

type RecipeChatHistory struct {
	gorm.Model
	// RecipeID     uint           `gorm:"uniqueIndex;"`
	// MessagesJSON pq.StringArray `gorm:"type:text[]"`
	Messages []RecipeChatHistoryMessage `gorm:"foreignKey:RecipeChatHistoryID"`
}

type RecipeChatHistoryMessage struct { // this is a single message, it's serialized and appended to the messages array
	gorm.Model
	RecipeChatHistoryID uint // Foreign key (belongs to RecipeChatHistory)
	UserPrompt          string
	GeneratedText       FunctionCallArgument `gorm:"type:json"` // recipeManger is serialized and placed here
}

// generated recipe json is given back as a json string and userInput is already provided as userPrompt(change the name of this variable)

// messages would be a serialized RecipeChatMessage; userInput, followed by generatedText, followed by userInput, etc.

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
// 		ImagePrompt string       `json:"image_prompt"`
// 	}{
// 		MainRecipe:  r.GeneratedRecipe.MainRecipe,
// 		SubRecipes:  r.GeneratedRecipe.SubRecipes,
// 		ImagePrompt: r.GeneratedRecipe.ImagePrompt,
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
