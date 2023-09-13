package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/models"
)

func getUserFromContext(c *gin.Context) (*models.User, error) {
	val, ok := c.Get("user")
	if !ok {
		return nil, errors.New("No user information")
	}

	user, ok := val.(*models.User)
	if !ok {
		return nil, errors.New("User information is of the wrong type")
	}

	return user, nil
}
