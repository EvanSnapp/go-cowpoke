package routes

//UpgradeStackRequest represents the data sent to upgrade a rancher stack
type UpgradeStackRequest struct {
	Catalog         string `json:"catalog" binding:"required"`
	Template        string `json:"template" binding:"required"`
	TemplateVersion string `json:"templateVersion" binding:"required"`
}

//StatusResponse represents the data returned by the status routes
type StatusResponse struct {
	Status string `json:"status"`
}
