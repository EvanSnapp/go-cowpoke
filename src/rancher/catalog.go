package rancher

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

//Template represents catalog template data retrived from the Rancher API
type Template struct {
	Name         string                 `json:"name"`
	ID           string                 `json:"id"`
	VersionLinks map[string]interface{} `json:"versionLinks"`
}

//TemplateVersion represents catalog template version data retrieved from the Rancher API
type TemplateVersion struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	CatalogID      string                 `json:"catalogId"`
	Version        string                 `json:"version"`
	DefaultVersion string                 `json:"defaultVersion"`
	TemplateFiles  map[string]interface{} `json:"files"`
}

//GetTemplateURL returns the URL associated with a catalog template at the specified version
func GetTemplateURL(catalog string, template string, version string) (*url.URL, error) {

	var data Template
	catalogID := fmt.Sprintf("%s:%s", url.PathEscape(catalog), url.PathEscape(template))
	catalogURL := os.Getenv("RANCHER_URL") + "/v1-catalog/templates/" + catalogID
	fmt.Println(catalogURL)

	if err := DoRequest(catalogURL, &data); err != nil {
		return nil, err
	}

	if val, found := data.VersionLinks[version]; found {
		//paranoia check...make sure that the found value is a string
		if templateVersionURL, isString := val.(string); isString {
			//putting on the tin foil hat and ensuring that the string is actually a URL
			if url, err := url.Parse(templateVersionURL); err == nil {
				return url, nil

			}
			return nil, fmt.Errorf("template version [%s] was found, but was not a valid URL", version)
		}

		return nil, errors.New("template version URL is not a string")
	}

	return nil, fmt.Errorf("could not find catalog template with version [%s]", version)
}

//GetTemplateVersion will retrieve the rancher and docker information for a catalog template
//at the specified version.
func GetTemplateVersion(u *url.URL) (*TemplateVersion, error) {
	//get the url of the template
	//get the data
	//return a data structure of template data
	return nil, nil
}
