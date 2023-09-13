package repository

import (
	"github.com/windoze95/culinaryai/internal/db"
	"github.com/windoze95/culinaryai/internal/models"
)

type UserRepository struct {
	UserDB *db.UserDB
}

func NewUserRepository(userDB *db.UserDB) *UserRepository {
	return &UserRepository{UserDB: userDB}
}

func (r *UserRepository) CreateUserAndSettings(user *models.User, settings *models.UserSettings) error {
	return r.UserDB.CreateUserAndSettings(user, settings)
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	return r.UserDB.GetUserByUsername(username)
}

// func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
// 	var user models.User
// 	if err := r.UserDB.GetUserByID(userID, &user); err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (r *UserRepository) UpdateUserSettingsOpenAIKey(userID uint, encryptedOpenAIKey string) error {
	return r.UserDB.UpdateUserSettingsOpenAIKey(userID, encryptedOpenAIKey)
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	return r.UserDB.UsernameExists(username)
}

func (r *UserRepository) PreloadUserByID(userID uint, user *models.User) error {
	return r.UserDB.PreloadUserByID(userID, user)
}
