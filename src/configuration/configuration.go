package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var requiredEnvVars = []string{"HOST_PORT", "RANCHER_URL", "RANCHER_USER_KEY", "RANCHER_USER_SECRET"}
var optionalEnvVars = []string{"SLACK_TOKEN", "SLACK_CHANNELS", "API_KEY"}

//Init initializes configuration and does basic validation. This should only be called once
func Init() []error {
	var errors []error

	//if we're running locally, crash if an .env file isn't found
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			panic("a .env file is required for local development")
		}
	}

	for _, optVar := range optionalEnvVars {
		if os.Getenv(optVar) == "" {
			log.Printf("Warning: configuration value '%s' is not set", optVar)
		}
	}

	for _, reqVar := range requiredEnvVars {
		if os.Getenv(reqVar) == "" {
			errors = append(errors, fmt.Errorf("'%s' is a required configuration value", reqVar))
		}
	}

	return errors
}

//GetSlackChannels returns a list of channels to use for slack notifications
func GetSlackChannels() []string {
	chans := []string{}
	return chans
}
