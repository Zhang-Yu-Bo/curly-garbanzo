package twitchAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func (u *UserInfo) Validation() {
	if u.DisplayName == "" {
		u.DisplayName = "none"
	}
	if u.Description == "" {
		u.Description = "none"
	}
	if u.ProfileImageURL == "" {
		u.ProfileImageURL = "https://static-cdn.jtvnw.net/jtv_user_pictures/b5272bba-0dbf-4d53-af96-a969755f9366-profile_image-300x300.png"
	}
}

type UserInfoList struct {
	User []UserInfo `json:"data"`
}

func GetUserInfoByName(loginAccount string) (UserInfo, error) {
	fmt.Printf("GetUserInfoByName: %s\n", loginAccount)
	var err error
	var req *http.Request

	url := "https://api.twitch.tv/helix/users?login=" + loginAccount
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return UserInfo{}, err
	}

	fmt.Println("token", GetAppAccessToken())
	fmt.Println("client id", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+GetAppAccessToken())
	req.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))

	var res *http.Response
	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		return UserInfo{}, err
	}

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
	defer res.Body.Close()

	if len(result.User) <= 0 {
		return UserInfo{}, errors.New("there is no user data in response")
	}

	return result.User[0], nil
}
