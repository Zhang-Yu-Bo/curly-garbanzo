package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/Zhang-Yu-Bo/curly-garbanzo/model/twitchAPI"
)

func GetHMACMessage(r *http.Request, requestBody []byte) string {
	return r.Header.Get(twitchAPI.TWITCH_MESSAGE_ID) +
		r.Header.Get(twitchAPI.TWITCH_MESSAGE_TIMESTAMP) +
		string(requestBody)
}

func GetHMAC(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
