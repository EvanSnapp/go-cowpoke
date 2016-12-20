package routes

import (
	"net/http"
	"rancher"

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

	if c.Bind(&data) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	environments, envErr := rancher.GetEnvironments()

	if envErr != nil {
		c.Error(envErr).SetType(gin.ErrorTypePublic)
		return
	}

	//get the template version the target stack(s) should be upgraded to
	tmplVersion, tmplErr := rancher.GetTemplateVersion(data.Catalog, data.Template, data.TemplateVersion)

	if tmplErr != nil {
		c.Error(tmplErr).SetType(gin.ErrorTypePublic)
		return
	}

	/*TODOs:
	-do each upgrade in go routines to avoid the O(n^2) loop
		--upgrading is a resource intensive process for Rancher (e.g. pulling images, networking, etc)
		--so a throttling mechanism would need to be implemented
	*/
	//attempt to upgrade stacks in all the environment the svc account has access to
	for _, env := range environments {
		stacks, _ := rancher.GetStacksToUpgrade(env, tmplVersion)

		for _, stack := range stacks {
			if err := rancher.UpgradeStack(stack, tmplVersion); err != nil {
				c.Error(err).SetType(gin.ErrorTypePublic)
			}
		}
	}

	//TODO:
	//if a stack was upgraded in at least one environment
	//send a 200 along and state which one's were good and bad
	if len(c.Errors) == 0 {
		c.JSON(200, "we made it! no send a better response!")
	}
}
