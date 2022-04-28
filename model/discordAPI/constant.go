package discordAPI

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var discordWebhook string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
		return
	}

	discordWebhook = os.Getenv("DISCORD_WEBHOOK")
	if discordWebhook == "" {
		fmt.Println("discord webhook url is empty.")
	}
}
