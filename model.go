package main

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
}

type UserSettings struct {
	gorm.Model
	UserID             uint `gorm:"index"`
	EncryptedOpenAIKey string
	// MFASecret string
}

type Recipe struct {
	gorm.Model
	Title              string
	Content            string
	Tags               []Tag `gorm:"many2many:recipe_tags;"`
	GeneratedBy        *User `gorm:"foreignKey:GeneratedByUserID"`
	GeneratedByUserID  uint
	GuidingContentID   uint
	GuidingContentUID  uuid.UUID
	GuidingContent     *GuidingContent `gorm:"foreignKey:GuidingContentID"`
	UserPrompt         string
	GenerationComplete bool
}

type GuidingContent struct {
	gorm.Model
	UserID uint `gorm:"index"` // Reference to the user who created the content
	UID    uuid.UUID
	// DietaryRestrictions string // Specific dietary restrictions
	SupportingResearch string // Supporting research to help convey the user's expectations
	Instructions       string // Additional instructions or guidelines
}

type Tag struct {
	gorm.Model
	Name string `gorm:"index"`
}
