package slack

import (
	"configuration"
	"fmt"
	"os"

	"errors"

	"github.com/nlopes/slack"
)

//Client provides convience functions for sending notifications through Slack
type Client struct {
	Client    *slack.Client
	Channels  []string
	AuthToken string
}

//SendToAll will send the provided message to all slack channels configured in the client
func (c Client) SendToAll(message string) []error {
	var errors []error

	for _, slackChan := range c.Channels {
		if _, _, err := c.Client.PostMessage(slackChan, message, slack.PostMessageParameters{}); err != nil {
			fmt.Println(err)
			errors = append(errors, err)
		}
	}

	return errors
}

//NewClient is a factory function for the Client struct
func NewClient() (Client, error) {

	token := os.Getenv("SLACK_TOKEN")
	channels := configuration.GetSlackChannels()

	if !(token != "" && len(channels) > 0) {
		fmt.Println("whoops")
		return Client{}, errors.New("no slack token and/or channels are configured")
	}

	slackClient := slack.New(token)
	client := Client{
		AuthToken: token,
		Channels:  channels,
		Client:    slackClient,
	}

	return client, nil
}
