package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Recipe is the model for a recipe.
type Recipe struct {
	gorm.Model
	Title              string
	Version            int            `gorm:"default:1"`
	Ingredients        Ingredients    `gorm:"type:jsonb"`
	Instructions       pq.StringArray `gorm:"type:text[]"`
	CookTime           int
	UnitSystem         UnitSystem     `gorm:"type:int"`
	LinkedRecipes      []*Recipe      `gorm:"many2many:recipe_linked_recipes;association_jointable_foreignkey:link_recipe_id"`
	LinkSuggestions    pq.StringArray `gorm:"type:text[]"`
	Hashtags           []*Tag         `gorm:"many2many:recipe_tags;"`
	ImagePrompt        string
	ImageURL           string
	CreatedByID        uint
	CreatedBy          *User `gorm:"foreignKey:CreatedByID"`
	PersonalizationUID uuid.UUID
	UserEdited         bool           `gorm:"default:false"`
	HistoryID          uint           `gorm:"unique;index"`
	History            *RecipeHistory `gorm:"foreignKey:HistoryID"`
	ForkedFromID       *uint
	ForkedFrom         *Recipe    `gorm:"foreignKey:ForkedFromID"`
	CreateType         RecipeType `gorm:"type:text"`
}

// RecipeHistory is the model for a recipe history.
type RecipeHistory struct {
	gorm.Model
	Messages []RecipeHistoryMessage `gorm:"foreignKey:RecipeHistoryID"`
}

// RecipeHistoryMessage is the model for a recipe history message.
type RecipeHistoryMessage struct {
	gorm.Model
	RecipeHistoryID uint // Foreign key (belongs to RecipeHistory)
	UserPrompt      string
	RecipeType      RecipeType `gorm:"type:text"`
	RecipeResponse  RecipeDef  `gorm:"type:jsonb"` // Embedded struct
}

// Tag is the model for a recipe hashtag.
type Tag struct {
	gorm.Model
	Hashtag string `gorm:"index:idx_hashtag;unique"`
}

// RecipeType is the type for the RecipeType enum.
type RecipeType string

// RecipeType enum values.
const (
	RecipeTypeChat            RecipeType = "chat"
	RecipeTypeBasedOn         RecipeType = "based_on"
	RecipeTypeRegenChat       RecipeType = "regen_chat"
	RecipeTypeCopycat         RecipeType = "copycat"
	RecipeTypeImportVision    RecipeType = "import_vision"
	RecipeTypeImportLink      RecipeType = "import_link"
	RecipeTypeImportCopypasta RecipeType = "import_text"
	RecipeTypeManualEntry     RecipeType = "user_input"
)
