package twitchAPI

import (
	"os"
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

func GetMessageSecret() string {
	return os.Getenv("MESSAGE_SECRET")
}
