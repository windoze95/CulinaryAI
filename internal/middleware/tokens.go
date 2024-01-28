package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/windoze95/saltybytes-api/internal/config"
)

// VerifyTokenMiddleware verifies the JWT token provided in the Authorization header.
func VerifyTokenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader // Token is directly provided in the Authorization header

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Env.JwtSecretKey.Value()), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Type assert to float64 (default for JSON numbers)
			if idFloat, ok := claims["user_id"].(float64); ok {
				// Convert to uint
				userID := uint(idFloat)
				// Set the userID in the context
				c.Set("user_id", userID)
				c.Next()
			} else {
				// Handle error: claim is not a float64
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id in token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
