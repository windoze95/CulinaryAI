package repository

import (
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/windoze95/saltybytes-api/internal/db"
	"github.com/windoze95/saltybytes-api/internal/models"
	"github.com/windoze95/saltybytes-api/internal/util"
)

type UserRepository struct {
	UserDB *db.UserDB
}

func NewUserRepository(userDB *db.UserDB) *UserRepository {
	return &UserRepository{UserDB: userDB}
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.UserDB.CreateUser(user); err != nil {
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

	user = util.StripSensitiveUserData(user)

	return user, nil
}

func (r *UserRepository) GetUserAuthByUsername(username string) (*models.User, error) {
	return r.UserDB.GetUserByUsername(username)
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user, err := r.UserDB.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	user = util.StripSensitiveUserData(user)

	return user, nil
}

func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	user, err := r.UserDB.GetPreloadedUserByID(userID)
	if err != nil {
		return nil, err
	}

	user = util.StripSensitiveUserData(user)

	return user, nil
}

func (r *UserRepository) GetPreloadedUserByID(userID uint) (*models.User, error) {
	user, err := r.UserDB.GetPreloadedUserByID(userID)
	if err != nil {
		return nil, err
	}

	user = util.StripSensitiveUserData(user)

	return user, nil
}

func (r *UserRepository) GetUserByFacebookID(facebookID string) (*models.User, error) {
	user, err := r.UserDB.GetUserByFacebookID(facebookID)
	if err != nil {
		return nil, err
	}

	user = util.StripSensitiveUserData(user)

	return user, nil
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
