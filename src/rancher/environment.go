package rancher

import "os"

//Environment represents the "projects" resource from the Rancher API
//the name reflects what it's called in the UI and in the future API v2
type Environment struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	State string            `json:"state"`
	Links map[string]string `json:"links"`
}

//GetEnvironments retrieves the environments that the
//service account has access to
func GetEnvironments() ([]*Environment, error) {
	var err error
	endpointURL := os.Getenv("RANCHER_URL") + "/v1/projects/"
	response := struct {
		Data []*Environment `json:"data"`
	}{}

	if e := DoRequest(endpointURL, &response); err != nil {
		err = e
	}

	//if no environments were found, it's an authorization issue not a 404 issue
	if len(response.Data) == 0 {
		err = ErrForbidden
	}

	return response.Data, err
}
