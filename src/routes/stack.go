package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddStackRoutes wires up all of the HTTP routes
//for the stack endpoint
func AddStackRoutes(r *gin.RouterGroup) {
	r.PATCH("/stack", UpgradeStack)
}

//UpgradeStack attempts to upgrade a Rancher stack in all authorized environments
func UpgradeStack(c *gin.Context) {
	var data UpgradeStackRequest
	statusCode := http.StatusOK
	msg := gin.H{"msg": "stack upgraded successfully"}

	if c.Bind(&data) != nil {
		statusCode = http.StatusBadRequest
		msg["msg"] = c.Errors
	}

	c.JSON(statusCode, msg)
}
