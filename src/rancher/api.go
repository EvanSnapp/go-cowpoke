package rancher

import (
	"encoding/json"
	"net/http"
	"os"
)

var client = &http.Client{}

//DoRequest will make a call to the Rancher API and returns a map
//representing the JSON response
//TODO: do input validation (url parsing)
//can't think of a better way to do generic unmarshalling atm
func DoRequest(url string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("RANCHER_USER_KEY"), os.Getenv("RANCHER_USER_SECRET"))

	res, err := client.Do(req)

	//TODO: probably need different error types
	//ex: req bad or decoder error = http 500
	//otherwise unauthorized, etc
	//TODO: handle unauthorized responses
	if err != nil {
		return err
	}

	json.NewDecoder(res.Body).Decode(&response)
	return nil
}
