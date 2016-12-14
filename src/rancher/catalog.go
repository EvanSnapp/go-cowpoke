package rancher

import (
	"fmt"
	"net/url"
	"os"
	"rancher/types"
	"reflect"
)

//GetTemplateURL returns the URL associated with a catalog template at the specified version
func GetTemplateURL(catalog string, template string, version string) (string, error) {

	var data *types.Template
	var err error
	catalogID := fmt.Sprintf("%s:%s", url.PathEscape(catalog), url.PathEscape(template))
	catalogURL := os.Getenv("RANCHER_URL") + "/v1-catalog/templates/" + catalogID

	if e := DoRequest(catalogURL, &data); err != nil {
		err = e
	}

	templateURL, found := data.VersionLinks[version]

	if !found {
		err = ErrNotFound
	}

	return templateURL, ErrNotFound
}

//GetTemplateVersion will retrieve the rancher and docker information for a catalog template
//at the specified version.
func GetTemplateVersion(catalog string, template string, version string) (*types.TemplateVersion, error) {
	var data *types.TemplateVersion
	var err error

	templateURL, e := GetTemplateURL(catalog, template, version)

	if err != nil {
		err = e
	}

	if e2 := DoRequest(templateURL, &data); e2 != nil {
		err = e2
	}

	if (data == nil || reflect.DeepEqual(data, types.TemplateVersion{})) {
		err = ErrNotFound
	}

	return data, err
}
