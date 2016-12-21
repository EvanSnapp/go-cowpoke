package routes

import (
	"fmt"
	"net/http"
	"rancher"
	"rancher/types"
	"slack"

	"github.com/gin-gonic/gin"
)

//AddStackRoutes wires up all of the HTTP routes
//for the stack endpoint
func AddStackRoutes(r *gin.RouterGroup) {
	r.PATCH("/stack", UpgradeStack)
}

//UpgradeStack attempts to upgrade a Rancher stack in all authorized environments
func UpgradeStack(c *gin.Context) {
	var requestData UpgradeStackRequest

	if c.Bind(&requestData) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	environments, envErr := rancher.GetEnvironments()

	if envErr != nil {
		c.Error(envErr).SetType(gin.ErrorTypePublic)
		return
	}

	//get the template version the target stack(s) should be upgraded to
	tmplVersion, tmplErr := rancher.GetTemplateVersion(requestData.Catalog, requestData.Template, requestData.TemplateVersion)

	if tmplErr != nil {
		c.Error(tmplErr).SetType(gin.ErrorTypePublic)
		return
	}

	/*TODOs:
	-refactor this code into better abstractions
	-do each upgrade in go routines to avoid the O(n^2) loop
		--upgrading is a resource intensive process for Rancher (e.g. pulling images, networking, etc)
		--so a throttling mechanism would need to be implemented
	-better HTTP code and responses
	*/
	//attempt to upgrade stacks in all the environment the svc account has access to
	results := []types.StackUpgradeResult{}
	responseCode := 200 //shourtcut for determining HTTP code
	startMsgFmt := "upgrading stack [%s] in environment [%s] from [%s] to [%s:%s]"
	successMsg := "successfully upgraded stack [%s] in environment [%s] from [%s] to [%s:%s]"
	failMsg := "failed to upgrade stack [%s] in environment [%s] from [%s] to [%s:%s] \n Error: %s"

	for _, env := range environments {
		stacks, _ := rancher.GetStacksToUpgrade(env, tmplVersion)

		for _, stack := range stacks {

			startMsg := fmt.Sprintf(startMsgFmt, stack.Name, env.Name, stack.ExternalID, requestData.Template, requestData.TemplateVersion)
			endMsg := fmt.Sprintf(successMsg, stack.Name, env.Name, stack.ExternalID, requestData.Template, requestData.TemplateVersion)

			go doNotification(startMsg)

			result := rancher.UpgradeStack(stack, tmplVersion)
			results = append(results, result)

			if result.UpgradedTo == "" {
				responseCode = 500
				endMsg = fmt.Sprintf(failMsg, stack.Name, env.Name, stack.ExternalID, requestData.Template, requestData.TemplateVersion, result.Error)
			}

			go doNotification(endMsg)
		}
	}

	//TODO: this response data needs to be improved
	c.JSON(responseCode, gin.H{"msg": "results from upgrading stack(s)", "results": results})
}

//helper function for sending slack notifcations
func doNotification(msg string) {
	if client, err := slack.NewClient(); err == nil {
		client.SendToAll(msg)
	}
}
