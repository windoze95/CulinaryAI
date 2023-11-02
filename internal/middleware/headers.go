package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckIDHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logging the entire header
		for key, value := range c.Request.Header {
			log.Printf("%s: %s\n", key, strings.Join(value, ", "))
		}

		// Logging parts of the request
		log.Println("Method:", c.Request.Method)
		log.Println("URL:", c.Request.URL.String())
		log.Println("Protocol:", c.Request.Proto)

		idHeaderValue := c.GetHeader("X-SaltyBytes-Identifier")
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
