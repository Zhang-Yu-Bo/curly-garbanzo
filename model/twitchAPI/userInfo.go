package twitchAPI

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UserInfo struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	CreatedAt       string `json:"createed_at"`
}

type UserInfoList struct {
	User []UserInfo `json:"data"`
}

func GetUserInfoByName(loginAccount, token string) (UserInfo, error) {
	if token == "" {
		token = appAccessToken
	}

	url := "https://api.twitch.tv/helix/users?login=" + loginAccount

	var err error
	var req *http.Request
	client := &http.Client{}
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return UserInfo{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Client-Id", clientID)

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return UserInfo{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		return UserInfo{}, errors.New("request was invalid")
	}
	if res.StatusCode == http.StatusUnauthorized {
		return UserInfo{}, errors.New("authorization failed")
	}

	result := UserInfoList{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return UserInfo{}, err
	}

	if len(result.User) > 0 {
		return UserInfo{}, errors.New("there is no user data in response")
	}

	return result.User[0], nil
}
