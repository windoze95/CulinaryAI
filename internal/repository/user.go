package repository

import (
	"github.com/windoze95/culinaryai/internal/db"
	"github.com/windoze95/culinaryai/internal/models"
)

type UserRepository struct {
	userDB *db.UserDB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.userDB.CreateUser(user)
}

func (r *UserRepository) PreloadUserByID(userID uint, user *models.User) error {
	return r.userDB.PreloadUserByID(userID, user)
}
