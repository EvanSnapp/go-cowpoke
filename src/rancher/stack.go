package rancher

import (
	"fmt"
	"strconv"
	"strings"
)

//Stack represents the "environments" resource from the Rancher API
//the name reflects what it's called in the UI and in the future API v2
type Stack struct {
	ID                   string            `json:"id"`
	ExternalID           string            `json:"externalId"`
	Name                 string            `json:"name"`
	RancherEnvironmentID string            `json:"accountID"`
	DockerCompseYML      string            `json:"dockerCompose"`
	RancherComposeYML    string            `json:"rancherCompose"`
	EnvironmentVars      map[string]string `json:"environment"`
	ActionURLs           map[string]string `json:"actions"`
}

//IsUpgradableTo determines if a stack can be upgraded to a particular
//catalog template version id
//TODO: write tests for this
func (s *Stack) IsUpgradableTo(tmplVerID string) bool {

	/*
			   Rules:
			     - the template version id must be formatted correctly (i.e. <catalog>:<template>:<version num>)
			     - the stack has to be created from a catalog template (i.e catalog://<template version id>)
		       - the stack has to have created from the same catalog and template as the provided template version id
		       - the version number of the template the stack was created from must be lower than the number in the provided id
	*/
	var extIDParts []string
	isUpgradeable := false
	tmplVerIDParts := strings.Split(tmplVerID, ":")

	if len(tmplVerIDParts) == 3 {
		extIDParts = strings.Split(s.ExternalID, "//")

		if len(extIDParts) == 2 {
			extIDParts = strings.Split(extIDParts[1], ":")
			extIDVerNum, parseErr1 := strconv.Atoi(extIDParts[2])
			tmplVerIDNum, parseErr2 := strconv.Atoi(tmplVerIDParts[2])

			if (parseErr1 == nil && parseErr2 == nil) &&
				len(extIDParts) == 3 &&
				(extIDParts[0] == tmplVerIDParts[0] && extIDParts[1] == tmplVerIDParts[1] && extIDVerNum < tmplVerIDNum) {
				isUpgradeable = true
			}
		}
	}

	return isUpgradeable
}

//GetStacksToUpgrade retrieves all of the stacks in an environment that were created
//from the specified catalog template
func GetStacksToUpgrade(env *Environment, templateVersion *TemplateVersion) ([]*Stack, error) {
	var err error
	var stacks []*Stack
	endpointURL := env.Links["environments"]
	response := struct {
		Data []*Stack `json:"data"`
	}{}

	if e := DoRequest(endpointURL, &response); e != nil {
		err = e
	}

	if len(response.Data) == 0 {
		err = ErrNotFound
	}

	for _, stack := range response.Data {
		if stack.IsUpgradableTo(templateVersion.ID) {
			fmt.Println("found a stack to upgrade!")
			stacks = append(stacks, stack)
		}
	}

	if len(stacks) == 0 {
		err = ErrNotFound
	}

	return stacks, err
}
