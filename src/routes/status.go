package routes

import (
	"github.com/gin-gonic/gin"
)

//AddStatusRoutes wires up all of the HTTP routes
//for the status endpoint
func AddStatusRoutes(api *gin.RouterGroup) {
	api.GET("/_status", GetStatus)
}

//GetStatus is a smoke test endpoint for health checks
func GetStatus(c *gin.Context) {
	c.String(200, "service up")
}
