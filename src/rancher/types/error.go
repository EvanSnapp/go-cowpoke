package types

//APIError represents and error related to calling Rancher API endpoints
//the error can come from the API itself or by inspecting payloads, etc
type APIError struct {
	error
	URL     string
	Status  int    `json:"status,string"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}
