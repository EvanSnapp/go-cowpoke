package types

//Environment represents the "projects" resource from the Rancher API
//the name reflects what it's called in the UI and in the future API v2
type Environment struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	State string            `json:"state"`
	Links map[string]string `json:"links"`
}
