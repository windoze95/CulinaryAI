package util

import (
	"github.com/windoze95/saltybytes-api/internal/models"
)

func StripSensitiveUserData(user *models.User) *models.User {
	// Strip sensitive information from the user object
	user.Auth = models.UserAuth{}
	user.Settings.EncryptedOpenAIKey = ""
	return user
}
