package middleware

import (
	"net/http"
	"rancher/types"

	"github.com/gin-gonic/gin"
)

//Errors is a middleware responsible for sending HTTP responses when errors occur
func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		//if the request has been aborted then request should not be modified further
		if !c.IsAborted() {
			code := http.StatusInternalServerError
			msg := "an unexpected error occurred"
			//By convention only get the last "public" error that
			//was added to the req context. this represents errors that happen
			//at the route level
			if err := c.Errors.ByType(gin.ErrorTypePublic).Last(); err != nil {
				innerErr := err.Err
				//if the error was a rancher api specific error use it's status and message
				if rancherErr, ok := innerErr.(types.APIError); ok {
					code = rancherErr.Status
					msg = rancherErr.Message
				}

				c.JSON(code, gin.H{"msg": msg})
			}
			c.Abort()
		}
	}
}
