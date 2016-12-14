package rancher

import (
	"fmt"
	"rancher/types"
)

//GetStacksToUpgrade retrieves all of the stacks in an environment that were created
//from the specified catalog template
func GetStacksToUpgrade(env *types.Environment, templateVersion *types.TemplateVersion) ([]*types.Stack, error) {
	var err error
	var stacks []*types.Stack
	endpointURL := env.Links["environments"]
	response := struct {
		Data []*types.Stack `json:"data"`
	}{}

	if e := DoRequest(endpointURL, &response); e != nil {
		err = e
	}

	if len(response.Data) == 0 {
		err = types.APIError{Status: 404, Message: "No stacks found"}
	}

	for _, stack := range response.Data {
		if stack.IsUpgradableTo(templateVersion.ID) {
			fmt.Println("found a stack to upgrade!")
			stacks = append(stacks, stack)
		}
	}

	if len(stacks) == 0 {
		err = types.APIError{Status: 404, Message: "No stacks found"}
	}

	return stacks, err
}
