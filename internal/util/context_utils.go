package util

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/models"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
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
