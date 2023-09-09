package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username         string
	Email            string
	HashedPassword   string
	Settings         UserSettings `gorm:"foreignKey:UserID"`
	CollectedRecipes []Recipe     `gorm:"many2many:user_collected_recipes;"`
	GuidingContent   GuidingContent
}

type UserSettings struct {
	gorm.Model
	UserID             uint `gorm:"index"`
	EncryptedOpenAIKey string
}

type GuidingContent struct {
	gorm.Model
	UserID       uint `gorm:"index"` // Reference to the user who created the content
	UID          uuid.UUID
	Requirements string // Additional instructions or guidelines
	// DietaryRestrictions string // Specific dietary restrictions
	// SupportingResearch string // Supporting research to help convey the user's expectations
}
