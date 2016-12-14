package rancher

import (
	"os"
	"rancher/types"
)

//GetEnvironments retrieves the environments that the
//service account has access to
func GetEnvironments() ([]*types.Environment, error) {
	var err error
	endpointURL := os.Getenv("RANCHER_URL") + "/v1/projects/"
	response := struct {
		Data []*types.Environment `json:"data"`
	}{}

	if e := DoRequest(endpointURL, &response); err != nil {
		err = e
	}

	//if no environments were found, it's an authorization issue not a 404 issue
	if len(response.Data) == 0 {
		err = types.APIError{Status: 500, Message: "No environments could be found"}
	}

	return response.Data, err
}
