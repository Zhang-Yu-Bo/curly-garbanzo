package twitchAPI

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
)

var appAccessToken = os.Getenv("APP_ACCESS_TOKEN")

func GetAppAccessToken() string {
	if checkTokenValid(appAccessToken) {
		return appAccessToken
	}
	if err := updateToken(); err != nil {
		fmt.Printf("updateToken failed: token = %s, reason = %s\n", appAccessToken, err)
		return ""
	}
	os.Setenv("APP_ACCESS_TOKEN", appAccessToken)
	return appAccessToken
}

func checkTokenValid(token string) bool {
	url := "https://id.twitch.tv/oauth2/validate"

	var err error
	var req *http.Request
	client := &http.Client{}
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return false
	}
	req.Header.Add("Authorization", "OAuth "+token)

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		fmt.Printf("checkTokenValid: token = %s, reason = %s\n", token, err)
		return false
	}
	if res.StatusCode == http.StatusOK {
		return true
	}
	if res.StatusCode == http.StatusUnauthorized {
		return false
	}
	return false
}

func updateToken() error {
	url := "https://id.twitch.tv/oauth2/token"

	payload := strings.NewReader("client_id=" + os.Getenv("TWITCH_CLIENT_ID") +
		"&client_secret=" + os.Getenv("TWITCH_CLIENT_SECRET") +
		"&grant_type=client_credentials",
	)

	var err error
	var req *http.Request
	client := &http.Client{}

	if req, err = http.NewRequest("POST", url, payload); err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return err
	}

	var js *simplejson.Json
	if js, err = simplejson.NewFromReader(res.Body); err != nil {
		return err
	}
	defer res.Body.Close()

	var tempToken string
	if tempToken, err = js.Get("access_token").String(); err != nil {
		return err
	}
	appAccessToken = tempToken
	return nil
}
