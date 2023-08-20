package main

import (
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
	Title             string
	Content           string
	GeneratedBy       *User `gorm:"foreignKey:GeneratedByUserID"`
	GeneratedByUserID uint
	Tags              []Tag `gorm:"many2many:recipe_tags;"`
	UserPrompt        string
	Status            string
}

type GuidingContent struct {
	gorm.Model
	UserID uint // Reference to the user who created the content
	// DietaryRestrictions string // Specific dietary restrictions
	SupportingResearch string // Supporting research to help convey the user's expectations
	Instructions       string // Additional instructions or guidelines
}

type Tag struct {
	gorm.Model
	Name string `gorm:"index"`
}
