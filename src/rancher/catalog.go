package rancher

import (
	"fmt"
	"net/url"
	"os"
	"rancher/types"
)

//GetTemplateURL returns the URL associated with a catalog template at the specified version
func GetTemplateURL(catalog string, template string, version string) (*url.URL, error) {

	var data types.Template
	catalogID := fmt.Sprintf("%s:%s", url.PathEscape(catalog), url.PathEscape(template))
	catalogURL := os.Getenv("RANCHER_URL") + "/v1-catalog/templates/" + catalogID

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
			return nil, ErrServer
		}

		return nil, ErrServer
	}

	return nil, ErrNotFound
}

//GetTemplateVersion will retrieve the rancher and docker information for a catalog template
//at the specified version.
func GetTemplateVersion(catalog string, template string, version string) (*types.TemplateVersion, error) {
	var data *types.TemplateVersion

	templateURL, err := GetTemplateURL(catalog, template, version)

	if err != nil {
		return nil, err
	}

	if err := DoRequest(templateURL.String(), &data); err != nil {
		return nil, err
	}

	return data, nil
}
