package rancher

import (
	"fmt"
	"net/url"
	"os"
	"rancher/types"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

//GetTemplateURL returns the URL associated with a catalog template at the specified version
func GetTemplateURL(catalog string, template string, version string) (string, error) {

	var data *types.Template

	catalogID := fmt.Sprintf("%s:%s", url.PathEscape(catalog), url.PathEscape(template))
	catalogURL := os.Getenv("RANCHER_URL") + "/v1-catalog/templates/" + catalogID

	if e := DoRequest(catalogURL, &data); e != nil {
		return "", e
	}

	templateURL, found := data.VersionLinks[version]

	if !found {
		return "", types.APIError{Status: 404, Message: "url not found for template"}
	}

	return templateURL, nil
}

//GetTemplateVersion will retrieve the rancher and docker information for a catalog template
//at the specified version.
func GetTemplateVersion(catalog string, template string, version string) (*types.TemplateVersion, error) {
	var data *types.TemplateVersion
	templateURL, e := GetTemplateURL(catalog, template, version)
	spew.Dump(e)

	if e != nil {
		return nil, e
	}

	if e2 := DoRequest(templateURL, &data); e2 != nil {
		return nil, e2
	}

	if (data == nil || reflect.DeepEqual(data, types.TemplateVersion{})) {
		return nil, types.APIError{
			Status:  404,
			Message: "could not get template version from blah blah",
		}
	}

	return data, nil
}
