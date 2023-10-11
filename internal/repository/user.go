package repository

import (
	"github.com/windoze95/saltybytes-api/internal/db"
	"github.com/windoze95/saltybytes-api/internal/models"
)

type UserRepository struct {
	UserDB *db.UserDB
}

func NewUserRepository(userDB *db.UserDB) *UserRepository {
	return &UserRepository{UserDB: userDB}
}

func (r *UserRepository) CreateUser(user *models.User, settings *models.UserSettings, gc *models.GuidingContent) error {
	return r.UserDB.CreateUser(user, settings, gc)
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	return r.UserDB.GetUserByUsername(username)
}

func (r *UserRepository) GetUserByFacebookID(facebookID string) (*models.User, error) {
	return r.UserDB.GetUserByFacebookID(facebookID)
}

func (r *UserRepository) UpdateUserEmail(userID uint, email string) error {
	return r.UserDB.UpdateUserEmail(userID, email)
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

func (r *UserRepository) UpdateGuidingContent(userID uint, updatedGC *models.GuidingContent) error {
	return r.UserDB.UpdateGuidingContent(userID, updatedGC)
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	return r.UserDB.UsernameExists(username)
}

func (r *UserRepository) GetPreloadedUserByID(userID uint) (*models.User, error) {
	return r.UserDB.GetPreloadedUserByID(userID)
}
