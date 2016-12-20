package types

//StackUpgradeResult contains data relating to upgrading
//a stack in a particular environment
type StackUpgradeResult struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	UpgradedTo  string `json:"upgradedTo"`
	Error       string `json:"error"`
}
