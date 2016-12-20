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
func UpgradeStack(s *types.Stack, v *types.TemplateVersion) error {
	var curStackState *types.Stack
	upgradeURL, canUpgrade := s.ActionURLs["upgrade"]

	//something could have happened between getting the stack information and making this call so do
	//one last check before making the upgrade calls
	if !canUpgrade {
		msgFmt := "stack '%s'(id: %s) in environment %s can not be upgraded at this time"
		return types.APIError{Status: 422, Message: fmt.Sprintf(msgFmt, s.Name, s.ID, s.RancherEnvironmentID)}
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
		return upgradeErr
	}

	/*
		TODO:
			The original implementation is to poll until the stack is not in an "upgrading" state
			This doesn't scale. Here are some improvements that could be made:
				-if there is an issue with the data (e.g. invalid YAML) then the stack stays in the upgrading state forever
				-an exponential backoff instead of a static sleep between polls
				-a "circuit breaker" mechanism that stops polling after a certain amount of attempts
	*/
	for curStackState.State != "upgraded" {
		if upgradePollErr := DoRequest(curStackState.Links["self"], &curStackState); upgradePollErr != nil {
			return upgradePollErr
		}

		time.Sleep(500 * time.Millisecond)
	}

	//if we get to this point then the upgrade was successful and the upgrade should be finished
	//TODO: return an error if the call to finish failed
	if finishUpgradeErr := DoPost(curStackState.ActionURLs["finishupgrade"], "", ""); finishUpgradeErr != nil {
		return finishUpgradeErr
	}

	return nil
}
