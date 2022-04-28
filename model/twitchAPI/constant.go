package twitchAPI

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var appAccessToken string

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

	appAccessToken = os.Getenv("APP_ACCESS_TOKEN")
	if appAccessToken == "" {
		if err := updateToken(); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetMessageSecret() string {
	return os.Getenv("MESSAGE_SECRET")
}
