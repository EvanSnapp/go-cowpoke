package configuration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

//RequiredEnvVars is  a list of the env vars required for the service to function
var RequiredEnvVars = []string{"HOST_PORT", "RANCHER_URL", "RANCHER_USER_KEY", "RANCHER_USER_SECRET"}

//OptionalEnvVars is a list of the optional env vars for the service
var OptionalEnvVars = []string{"SLACK_TOKEN", "SLACK_CHANNELS", "API_KEY"}

//Init initializes configuration and does basic validation. This should only be called once
func Init() []error {
	var errors []error

	//if we're running locally, crash if an .env file isn't found
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			panic("a .env file is required for local development")
		}
	}

	for _, optVar := range OptionalEnvVars {
		if os.Getenv(optVar) == "" {
			log.Printf("Warning: configuration value '%s' is not set", optVar)
		}
	}

	for _, reqVar := range RequiredEnvVars {
		if os.Getenv(reqVar) == "" {
			errors = append(errors, fmt.Errorf("'%s' is a required configuration value", reqVar))
		}
	}

	return errors
}

//GetSlackChannels returns a list of channels to use for slack notifications
func GetSlackChannels() []string {
	var chans []string

	if chansVar := os.Getenv("SLACK_CHANNELS"); chansVar != "" {
		for _, chanVar := range strings.Split(chansVar, ",") {
			if trimmedChanVar := strings.TrimSpace(chanVar); trimmedChanVar != "" {
				chans = append(chans, trimmedChanVar)
			}
		}
	}

	return chans
}
