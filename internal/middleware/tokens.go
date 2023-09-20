package middleware

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/windoze95/culinaryai/internal/config"
)

func VerifyTokenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("VerifyTokenMiddleware", c.Request.Cookies())
		log.Println("VerifyTokenMiddleware", c.Request.Header.Get("Cookie"))
		log.Println("VerifyTokenMiddleware", c.Request.Header.Get("Authorization"))
		log.Println("VerifyTokenMiddleware", c.Request.Header.Get("auth_token"))
		cookie, err := c.Cookie("auth_token") // Fetch auth_token cookie
		if err != nil {
			log.Println("VerifyTokenMiddleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		log.Println("cookie", cookie)

		tokenString := cookie // Token is fetched from the cookie

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Env.JwtSecretKey.Value()), nil
		})
		if err != nil {
			log.Println("VerifyTokenMiddleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		log.Println("token", token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("claims", claims)
			log.Println("claims[user_id]", claims["user_id"])
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			log.Println("VerifyTokenMiddleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
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
