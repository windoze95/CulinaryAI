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
	UnitSystem   int    `gorm:"default:1"` // 1 = US Customary, 2 = Metric
	Requirements string // Additional instructions or guidelines
	// DietaryRestrictions string // Specific dietary restrictions
	// SupportingResearch string // Supporting research to help convey the user's expectations
}

func (gc *GuidingContent) BeforeCreate(tx *gorm.DB) (err error) {
	if gc.UnitSystem != 1 && gc.UnitSystem != 2 {
		gc.UnitSystem = 1 // Default to 1 if not 1 or 2
	}
	return nil
}

func (gc *GuidingContent) BeforeUpdate(tx *gorm.DB) (err error) {
	if gc.UnitSystem != 1 && gc.UnitSystem != 2 {
		gc.UnitSystem = 1 // Default to 1 if not 1 or 2
	}
	return nil
}
