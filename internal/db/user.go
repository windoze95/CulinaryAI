package db

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/windoze95/culinaryai/internal/models"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(gormDB *gorm.DB) *UserDB {
	return &UserDB{DB: gormDB}
}

func (db *UserDB) CreateUserAndSettings(user *models.User, settings *models.UserSettings) error {
	tx := db.DB.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	settings.UserID = user.ID
	if err := tx.Create(settings).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (db *UserDB) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (db *UserDB) GetUserByID(userID uint, user *models.User) error {
// 	return db.DB.Preload("Settings").Where("id = ?", userID).First(user).Error
// }

func (db *UserDB) UpdateUserSettingsOpenAIKey(userID uint, encryptedOpenAIKey string) error {
	return db.DB.Model(&models.UserSettings{}).Where("user_id = ?", userID).Update("EncryptedOpenAIKey", encryptedOpenAIKey).Error
}

func (db *UserDB) UsernameExists(username string) (bool, error) {
	lowercaseUsername := strings.ToLower(username)
	var user models.User
	if err := db.DB.Where("LOWER(username) = ?", lowercaseUsername).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *UserDB) PreloadUserByID(userID uint, user *models.User) error {
	return db.DB.Preload("Settings").Preload("GuidingContent").Where("id = ?", userID).First(user).Error
}
