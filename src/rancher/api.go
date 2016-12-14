package rancher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"rancher/types"

	"github.com/davecgh/go-spew/spew"
)

var client = &http.Client{}

//DoRequest will make a call to the Rancher API and decodes
//the JSON response into the supplied interface
func DoRequest(uri string, response interface{}) error {
	fmt.Println("doing call for " + uri)
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
		hai := types.APIError{}
		json.NewDecoder(res.Body).Decode(&hai)
		hai.URL = uri
		spew.Dump(hai)
		fmt.Println(hai.Status)
		fmt.Println(hai.Message)
		return hai
	}

	json.NewDecoder(res.Body).Decode(&response)

	return nil
}
