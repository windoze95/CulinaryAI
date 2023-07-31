package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username         string
	HashedPassword   string
	Settings         UserSettings `gorm:"foreignKey:UserID"`
	CollectedRecipes []Recipe     `gorm:"many2many:user_recipes;"`
}

type UserSettings struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	OpenAIKey string
	// MFASecret string
}

type Recipe struct {
	gorm.Model
	Title       string
	Content     string
	Tags        []Tag `gorm:"many2many:recipe_tags;"`
	GeneratedBy *User `gorm:"foreignKey:UserID"`
	UserID      uint
	UserPrompt  *string
	DeletedAt   *time.Time
}

type Tag struct {
	gorm.Model
	Name string `gorm:"index"`
}
