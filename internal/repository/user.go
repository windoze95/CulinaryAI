package repository

import (
	"errors"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/windoze95/saltybytes-api/internal/models"
	"github.com/windoze95/saltybytes-api/internal/util"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	tx := r.DB.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		// Check for unique constraints
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			if strings.Contains(pgErr.Error(), "username") {
				return nil, errors.New("username already in use")
			} else if strings.Contains(pgErr.Error(), "email") {
				return nil, errors.New("email already in use")
			}
		}
		return nil, err
	}

	// user = util.StripSensitiveUserData(user)

	return util.StripSensitiveUserData(user), nil
}

func (r *UserRepository) GetUserAuthByUsername(username string) (*models.User, error) {
	// return r.UserDB.GetUserByUsername(username)
	var user models.User
	if err := r.DB.Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	// user, err := r.UserDB.GetUserByUsername(username)
	// if err != nil {
	// 	return nil, err
	// }
	var user models.User
	if err := r.DB.Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	// return &user, nil

	// user = util.StripSensitiveUserData(&user)

	return util.StripSensitiveUserData(&user), nil
}

func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	// user, err := r.UserDB.GetPreloadedUserByID(userID)
	// if err != nil {
	// 	return nil, err
	// }
	var user models.User
	if err := r.DB.Preload("Settings").
		Preload("GuidingContent").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	// user = *util.StripSensitiveUserData(&user)

	return util.StripSensitiveUserData(&user), nil
}

func (r *UserRepository) GetPreloadedUserByID(userID uint) (*models.User, error) {
	// user, err := r.UserDB.GetPreloadedUserByID(userID)
	// if err != nil {
	// 	return nil, err
	// }
	var user models.User
	if err := r.DB.Preload("Settings").
		Preload("GuidingContent").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	// user = util.StripSensitiveUserData(user)

	return util.StripSensitiveUserData(&user), nil
}

func (r *UserRepository) GetUserByFacebookID(facebookID string) (*models.User, error) {
	// user, err := r.UserDB.GetUserByFacebookID(facebookID)
	// if err != nil {
	// 	return nil, err
	// }

	var user models.User
	if err := r.DB.Where("facebook_id = ?", facebookID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return util.StripSensitiveUserData(&user), nil
}

func (r *UserRepository) UpdateUserEmail(userID uint, email string) error {
	err := r.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Email", email).Error
	if err != nil {
		log.Printf("Error updating user email: %v", err)
	}

	return err
}

// func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
// 	var user models.User
// 	if err := r.UserDB.GetUserByID(userID, &user); err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (r *UserRepository) UpdateUserSettingsOpenAIKey(userID uint, encryptedOpenAIKey string) error {
	// return r.UserDB.UpdateUserSettingsOpenAIKey(userID, encryptedOpenAIKey)
	err := r.DB.Model(&models.UserSettings{}).
		Where("user_id = ?", userID).
		Update("EncryptedOpenAIKey", encryptedOpenAIKey).Error
	if err != nil {
		log.Printf("Error updating user settings openai key: %v", err)
	}

	return err
}

func (r *UserRepository) UpdateGuidingContent(userID uint, updatedGC *models.GuidingContent) error {
	var existingGC models.GuidingContent

	// First, find the existing record
	err := r.DB.Where("user_id = ?", userID).
		First(&existingGC).Error
	if err != nil {
		log.Printf("Error retrieving existing guiding content: %v", err)
		return err
	}

	// Update fields
	existingGC.UnitSystem = updatedGC.UnitSystem
	existingGC.Requirements = updatedGC.Requirements

	// Perform the update
	err = r.DB.Save(&existingGC).Error
	if err != nil {
		log.Printf("Error saving updated guiding content: %v", err)
	}

	return err
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	lowercaseUsername := strings.ToLower(username)
	var user models.User
	err := r.DB.Where("LOWER(username) = ?", lowercaseUsername).
		First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
