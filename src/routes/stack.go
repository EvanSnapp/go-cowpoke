package routes

import (
	"fmt"
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
	-abstract this away from this route function
	-do whole upgrades process in seperate go routines to avoid the O(n^2) loop
	-will most likely need to throttle calls as upgrading a stack is a resource
	 intensive operation for rancher
	*/
	for _, env := range environments {
		fmt.Printf("getting stacks we can upgrade in [%s]\n", env.Name)
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
