package types

//Template represents catalog template data retrived from the Rancher API
type Template struct {
	Name         string            `json:"name"`
	ID           string            `json:"id"`
	VersionLinks map[string]string `json:"versionLinks"`
}

//TemplateVersion represents catalog template version data retrieved from the Rancher API
type TemplateVersion struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	CatalogID      string            `json:"catalogId"`
	Version        string            `json:"version"`
	DefaultVersion string            `json:"defaultVersion"`
	TemplateFiles  map[string]string `json:"files"`
}
