package repository

import (
	"errors"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/windoze95/saltybytes-api/internal/models"
)

// UserRepository is a repository for interacting with users.
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user.
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

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Settings").
		Preload("Personalization").
		Preload("Subscription").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserAuthByUsername retrieves a user's authentication information by their username.
func (r *UserRepository) GetUserAuthByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Auth").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserEmail updates a user's email address.
func (r *UserRepository) UpdateUserEmail(userID uint, email string) error {
	err := r.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Email", email).Error
	if err != nil {
		log.Printf("Error updating user email: %v", err)
	}

	return err
}

// UpdateUserSettingsKeepScreenAwake updates a user's KeepScreenAwake setting.
func (r *UserRepository) UpdateUserSettingsKeepScreenAwake(userID uint, keepScreenAwake bool) error {
	err := r.DB.Model(&models.UserSettings{}).
		Where("user_id = ?", userID).
		Update("KeepScreenAwake", keepScreenAwake).Error
	if err != nil {
		log.Printf("Error updating user settings: %v", err)
	}

	return err
}

// UpdatePersonalization updates a user's personalization settings.
func (r *UserRepository) UpdatePersonalization(userID uint, updatedPersonalization *models.Personalization) error {
	var existingPersonalization models.Personalization

	// First, find the existing record
	err := r.DB.Where("user_id = ?", userID).
		First(&existingPersonalization).Error
	if err != nil {
		log.Printf("Error retrieving existing personalization: %v", err)
		return err
	}

	// Update fields
	existingPersonalization.UnitSystem = updatedPersonalization.UnitSystem
	existingPersonalization.Requirements = updatedPersonalization.Requirements
	existingPersonalization.UID = updatedPersonalization.UID

	// Perform the update
	err = r.DB.Save(&existingPersonalization).Error
	if err != nil {
		log.Printf("Error saving updated personalization: %v", err)
	}

	return err
}

// UsernameExists checks if a username already exists.
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
