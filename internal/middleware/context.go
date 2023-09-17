package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/service"
)

func AttachUserToContext(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.Set("user", nil)
			c.Next()
			return
		}

		userID, ok := userIDValue.(float64) // jwt-go defaults to float64 for numerical claims
		if !ok || userID == 0 {
			c.Set("user", nil)
			c.Next()
			return
		}

		user, err := userService.GetPreloadedUserByID(uint(userID))
		if err != nil {
			c.Set("user", nil)
		} else {
			user.HashedPassword = "" // Remove password from user object
			c.Set("user", user)
		}
		c.Next()
	}
}
