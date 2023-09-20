package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/config"
	"github.com/windoze95/culinaryai/internal/util"
)

func VerifyTokenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token") // Fetch auth_token cookie
		if err != nil {
			util.ClearAuthTokenCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided", "forceLogout": true})
			c.Abort()
			return
		}

		tokenString := cookie // Token is fetched from the cookie

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Env.JwtSecretKey.Value()), nil
		})
		if err != nil {
			util.ClearAuthTokenCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "forceLogout": true})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			util.ClearAuthTokenCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "forceLogout": true})
			c.Abort()
			return
		}
	}
}

// func VerifyTokenMiddleware(cfg *config.Config) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		tokenString := authHeader // Token is directly provided in the Authorization header

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return cfg.Env.JwtSecretKey.Value(), nil
// 		})

// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			c.Set("user_id", claims["user_id"])
// 			c.Next()
// 		} else {
// 			c.JSON(401, gin.H{"message": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}
// 	}
// }
