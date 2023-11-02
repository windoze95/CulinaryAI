package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckIDHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		idHeaderValue := c.GetHeader("X-SaltyBytes-Identifier")
		log.Printf("idHeaderValue: %v", idHeaderValue)
		if idHeaderValue != "SByt3sIDToken" { // This isn't a security measure as much as it is a cors configuration
			// If the header is absent or the value is incorrect, reject the request
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Otherwise, proceed with the request
		c.Next()
	}
}
