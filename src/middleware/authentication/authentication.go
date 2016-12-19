package authentication

import (
	"os"

	"github.com/gin-gonic/gin"
)

//Authenticate is middleware that checks a
//client supplied token against a configured key from the environment
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		configuredAPIKey := os.Getenv("API_KEY")
		token := c.Request.Header.Get("bearer")

		if configuredAPIKey == "" || token == configuredAPIKey {
			c.Next()
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
			c.Abort()
		}
	}
}
