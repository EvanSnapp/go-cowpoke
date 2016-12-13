package routes

import (
	"fmt"
	"net/http"
	"rancher"

	"github.com/davecgh/go-spew/spew"
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

	/*
		TODO:
			- do the upgrade in each environment (async if possible)
			- send a payload stating what environments were uploaded
	*/
	//What should we do if there were no stacks to upgrade? send a 404?
	//TODO: this should be concurrent as well (done in the upgrade function)
	for _, env := range environments {
		fmt.Printf("getting stacks we can upgrade in [%s]\n", env.Name)
		stacks, _ := rancher.GetStacksToUpgrade(env, tmplVersion)
		spew.Dump(stacks)
	}

	c.JSON(200, "we made it")
}
