package twitchAPI

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	clientID       string
	clientSecret   string
	appAccessToken string
	messageSecret  string
)

const (
	TWITCH_MESSAGE_ID        = "twitch-eventsub-message-id"
	TWITCH_MESSAGE_TIMESTAMP = "twitch-eventsub-message-timestamp"
	TWITCH_MESSAGE_SIGNATURE = "twitch-eventsub-message-signature"
	MESSAGE_TYPE             = "twitch-eventsub-message-type"
)

// Notification message types
const (
	MESSAGE_TYPE_VERIFICATION = "webhook_callback_verification"
	MESSAGE_TYPE_NOTIFICATION = "notification"
	MESSAGE_TYPE_REVOCATION   = "revocation"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
		return
	}

	clientID = os.Getenv("TWITCH_CLIENT_ID")
	if clientID == "" {
		fmt.Println("client ID is empty.")
	}

	clientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	if clientSecret == "" {
		fmt.Println("client secret is empty.")
	}

	appAccessToken = os.Getenv("APP_ACCESS_TOKEN")
	if appAccessToken == "" {
		if err := updateToken(); err != nil {
			fmt.Println(err.Error())
		}
	}

	messageSecret = os.Getenv("MESSAGE_SECRET")
	if messageSecret == "" {
		fmt.Println("message secret is empty.")
	}
}

func GetMessageSecret() string {
	return messageSecret
}
