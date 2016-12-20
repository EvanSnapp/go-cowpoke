package rancher

//TODO: centralize net/http code

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"rancher/types"
)

var client = &http.Client{}

//DoRequest will make a call to the Rancher API and decodes
//the JSON response into the supplied interface
func DoRequest(uri string, response interface{}) error {
	if _, err := url.Parse(uri); err != nil {
		return err
	}

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("RANCHER_USER_KEY"), os.Getenv("RANCHER_USER_SECRET"))
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		apiErr := types.APIError{}
		json.NewDecoder(res.Body).Decode(&apiErr)
		apiErr.URL = uri
		return apiErr
	}

	json.NewDecoder(res.Body).Decode(&response)

	return nil
}

//DoPost makes POST calls to the Rancher API
func DoPost(uri string, data interface{}, response interface{}) error {
	if _, parseErr := url.Parse(uri); parseErr != nil {
		return parseErr
	}

	dataBytes, marshalErr := json.Marshal(data)

	if marshalErr != nil {
		return marshalErr
	}

	req, reqErr := http.NewRequest("POST", uri, bytes.NewBuffer(dataBytes))

	if reqErr != nil {
		return reqErr
	}

	req.SetBasicAuth(os.Getenv("RANCHER_USER_KEY"), os.Getenv("RANCHER_USER_SECRET"))
	req.Header.Set("Content-Type", "application/json")
	res, resErr := client.Do(req)

	if resErr != nil {
		return resErr
	}

	if !(res.StatusCode == 200 || res.StatusCode == 202) {
		apiErr := types.APIError{}
		json.NewDecoder(res.Body).Decode(&apiErr)
		apiErr.URL = uri
		return apiErr
	}

	json.NewDecoder(res.Body).Decode(&response)

	return nil
}
