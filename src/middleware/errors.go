package middleware

import (
	"net/http"
	"rancher"

	"github.com/gin-gonic/gin"
)

//Errors is middleware that checks a
//client supplied token against a configured key from the environment
func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		//if the request has been aborted then an error has already occurred and
		//no more processing should happen
		if !c.IsAborted() {

			//By convention only get the last "public" error that
			//was added to the req context. this represents errors that happen
			//at the route level
			if err := c.Errors.ByType(gin.ErrorTypePublic).Last(); err != nil {
				innerErr := err.Err
				code := http.StatusInternalServerError
				errMsg := innerErr.Error()

				switch innerErr {
				case rancher.ErrForbidden:
					code = http.StatusForbidden
				case rancher.ErrNotFound:
					code = http.StatusNotFound
					break
				case rancher.ErrUnauthorized:
					code = http.StatusUnauthorized
					break
				case rancher.ErrServer:
				default:
					errMsg = "an unexpected error occurred"
					break
				}

				c.JSON(code, gin.H{"msg": errMsg})
				//not sure if this is needed
				c.Abort()
			}
		}
	}
}
