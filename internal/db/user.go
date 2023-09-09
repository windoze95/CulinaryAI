package db

import (
	"github.com/jinzhu/gorm"
	"github.com/windoze95/culinaryai/internal/models"
)

type UserDB struct {
	DB *gorm.DB
}

func (db *UserDB) CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func (db *UserDB) PreloadUserByID(userID uint, user *models.User) error {
	return db.DB.Preload("Settings").Preload("GuidingContent").Where("id = ?", userID).First(user).Error
}
