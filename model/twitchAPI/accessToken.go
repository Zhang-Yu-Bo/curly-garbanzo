package twitchAPI

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

func GetAppAccessToken() string {
	if checkTokenValid(appAccessToken) {
		return appAccessToken
	}
	if err := updateToken(); err != nil {
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
	req.Header.Add("Authorization", "OAuth "+appAccessToken)

	var res *http.Response
	if res, err = client.Do(req); err != nil {
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

	defer res.Body.Close()

	result := map[string]string{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	var exist bool
	if appAccessToken, exist = result["access_token"]; !exist {
		return errors.New("update access token failed, there is no access_token field in response json")
	}
	return nil
}
