package util

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/models"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	val, ok := c.Get("user")
	if !ok {
		return nil, errors.New("no user information")
	}

	user, ok := val.(*models.User)
	if !ok {
		return nil, errors.New("user information is of the wrong type")
	}

	return user, nil
}

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	val, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("no user ID information")
	}

	userID, ok := val.(uint)
	log.Println("userID:", userID)
	log.Println("ok:", ok)
	if !ok {
		return 0, errors.New("user ID information is of the wrong type")
	}

	return userID, nil
}

// ClearAuthTokenCookie clears the authentication token cookie from the client's browser.
func ClearAuthTokenCookie(c *gin.Context) {
	// Clear the auth_token cookie
	c.SetCookie(
		"auth_token",         // Cookie name
		"",                   // Empty value to clear the cookie
		-1,                   // Max age < 0 to expire the cookie immediately
		"/",                  // Path
		".api.saltybytes.ai", // Domain, set with leading dot for subdomain compatibility
		true,                 // Secure
		true,                 // HTTP only
	)
}
