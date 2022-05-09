package twitchAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type ChannelInfo struct {
	BroadcasterId       string `json:"broadcaster_id"`
	BroadcasterLogin    string `json:"broadcaster_login"`
	BroadcasterName     string `json:"broadcaster_name"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	GameId              string `json:"game_id"`
	GameName            string `json:"game_name"`
	Title               string `json:"title"`
	Delay               int    `json:"delay"`
}

type ChannelInfoList struct {
	Channel []ChannelInfo `json:"data"`
}

func GetChannelInfoById(broadcasterId string) (mChannel ChannelInfo, err error) {
	fmt.Printf("GetChannelInfoById: %s\n", broadcasterId)

	url := "https://api.twitch.tv/helix/channels?broadcaster_id=" + broadcasterId

	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return mChannel, err
	}

	fmt.Println("token", GetAppAccessToken())
	fmt.Println("client id", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+GetAppAccessToken())
	req.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))

	var res *http.Response
	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		return mChannel, err
	}

	if res.StatusCode == http.StatusBadRequest {
		return mChannel, errors.New("missing query parameter")
	}
	if res.StatusCode == http.StatusUnauthorized {
		return mChannel, errors.New("authorization failed")
	}
	if res.StatusCode == http.StatusInternalServerError {
		return mChannel, errors.New("internal server error; failed to get channel information")
	}

	result := ChannelInfoList{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return mChannel, err
	}
	defer res.Body.Close()

	if len(result.Channel) <= 0 {
		return mChannel, errors.New("there is no channel data in response")
	}

	return result.Channel[0], nil
}
