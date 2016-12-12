package rancher

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

var (
	//ErrUnauthorized is an error for Rancher API 401
	ErrUnauthorized = errors.New("unauthorized")
	//ErrNotFound is an error for Rancher API 404 or for 200s with no data (aka no environments found)
	ErrNotFound = errors.New("no data found")
	//ErrServer is for all other erros coming back from the Rancher API
	ErrServer = errors.New("something bad happened with rancher")
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

	//treating default status codes as an error as we don't have any data
	switch res.StatusCode {
	case 200:
		json.NewDecoder(res.Body).Decode(&response)
		return nil
	case 401:
		return ErrUnauthorized
	case 404:
		return ErrNotFound
	case 500:
		return ErrServer
	default:
		return ErrServer
	}
}
