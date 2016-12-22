package rancher

import (
	"fmt"
	"rancher/types"
	"time"
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
			stacks = append(stacks, stack)
		}
	}

	if len(stacks) == 0 {
		err = types.APIError{Status: 404, Message: "No stacks found"}
	}

	return stacks, err
}

//UpgradeStack takes a stack and upgrades it to the specified template version
/* TODOs:
- better input validation
-making this take string args instead of pointers may be better in order to make concurrency a bit safer
*/
func UpgradeStack(s *types.Stack, v *types.TemplateVersion) types.StackUpgradeResult {
	var curStackState *types.Stack

	upgradeURL, canUpgrade := s.ActionURLs["upgrade"]
	result := types.StackUpgradeResult{
		Name:        s.Name,
		Environment: s.RancherEnvironmentID,
	}

	//something could have happened between getting the stack information and making this call so do
	//one last check before making the upgrade calls
	if !canUpgrade {
		result.Error = "the stack is not in an upgradeable state"
		return result
	}

	//TODO: make this it's own type
	data := struct {
		ExternalID     string            `json:"externalId"`
		DockerCompose  string            `json:"dockerCompose"`
		RancherCompose string            `json:"rancherCompose"`
		Environment    map[string]string `json:"environment"`
	}{
		ExternalID:     fmt.Sprintf("catalog://%s", v.ID),
		DockerCompose:  v.TemplateFiles["docker-compose.yml"],
		RancherCompose: v.TemplateFiles["rancher-compose.yml"],
		Environment:    s.EnvironmentVars,
	}

	if upgradeErr := DoPost(upgradeURL, data, &curStackState); upgradeErr != nil {
		result.Error = "an unexpected error occured during the upgrade request"
		return result
	}

	/*
		TODO:
			The original implementation is to poll until the stack is not in an "upgrading" state
			This doesn't scale. Here are some improvements that could be made:
				-if there is an issue with the data (e.g. invalid YAML) then the stack stays in the upgrading state forever
				-an exponential backoff instead of a static sleep between polls
				-a "circuit breaker" mechanism that stops polling after a certain amount of attempts
				-reverting the stack back to it's original state if an error condition is hit
	*/
	for curStackState.State != "upgraded" {
		if upgradePollErr := DoRequest(curStackState.Links["self"], &curStackState); upgradePollErr != nil {
			result.Error = "an unexpected error occurred while validating the upgrade request"
			return result
		}

		time.Sleep(500 * time.Millisecond)
	}

	//if we get to this point then the upgrade was successful and the upgrade should be finished
	result.UpgradedTo = v.Version
	if finishUpgradeErr := DoPost(curStackState.ActionURLs["finishupgrade"], "", ""); finishUpgradeErr != nil {
		result.Error = "the upgrade succeeded but could not be finished"
	}

	return result
}
