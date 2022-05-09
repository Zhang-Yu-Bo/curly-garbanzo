package controller

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/Zhang-Yu-Bo/curly-garbanzo/model/discordAPI"
	"github.com/Zhang-Yu-Bo/curly-garbanzo/model/twitchAPI"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func ShowUserInfo(w http.ResponseWriter, r *http.Request) {
	var err error
	var userInfo twitchAPI.UserInfo
	vars := mux.Vars(r)

	userInfo, err = twitchAPI.GetUserInfoByName(vars["username"])
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	fmt.Fprintf(w, "%v\n", userInfo)
}

func EventSub(w http.ResponseWriter, r *http.Request) {

	var err error
	var bodyBuf []byte
	if bodyBuf, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	var js *simplejson.Json
	if js, err = simplejson.NewJson(bodyBuf); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		fmt.Println(err)
		return
	}

	// 驗證訊息是否來自 twitch
	secret := twitchAPI.GetMessageSecret()
	message := GetHMACMessage(r, bodyBuf)
	hmac := "sha256=" + GetHMAC(secret, message)
	// 不太正確，應該使用 crypto.timingSafeEqual類的 equal比較安全，但方便起見這邊先直接比較
	if hmac != r.Header.Get(twitchAPI.TWITCH_MESSAGE_SIGNATURE) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "validation failed, hmac is not equal to TWITCH_MESSAGE_SIGNATURE")
		fmt.Println("secret =", secret)
		fmt.Println("message =", message)
		fmt.Println("mHMAC =", hmac)
		fmt.Println("tHMAC =", r.Header.Get(twitchAPI.TWITCH_MESSAGE_SIGNATURE))
		fmt.Println("validation failed, hmac is not equal to TWITCH_MESSAGE_SIGNATURE")
		return
	}

	switch r.Header.Get(twitchAPI.MESSAGE_TYPE) {
	case twitchAPI.MESSAGE_TYPE_NOTIFICATION:

		// get login account from event sub
		var twitchAccount string
		if twitchAccount, err = js.Get("event").Get("broadcaster_user_login").String(); err != nil {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println(err)
			return
		}
		// get user id from event sub
		var twitchBroadcasterUserId string
		if twitchBroadcasterUserId, err = js.Get("event").Get("broadcaster_user_id").String(); err != nil {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println(err)
			return
		}

		// get channel infomation from twitch API
		var channelInfo twitchAPI.ChannelInfo
		if channelInfo, err = twitchAPI.GetChannelInfoById(twitchBroadcasterUserId); err != nil {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println(err)
			return
		}
		// get user infomation from twitch API
		var userInfo twitchAPI.UserInfo
		if userInfo, err = twitchAPI.GetUserInfoByName(twitchAccount); err != nil {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println(err)
			return
		}

		userInfo.Validation()
		err = discordAPI.SendMessage(discordAPI.MessageOption{
			TagEveryone:       false,
			Content:           "** ﾚ(ﾟ∀ﾟ;)ﾍ 單兵注意 " + userInfo.DisplayName + " 開直播囉 ﾍ( ﾟ∀ﾟ;)ﾉ **",
			EmbedEnable:       true,
			EmbedTitle:        userInfo.DisplayName + " - " + channelInfo.Title,
			EmbedDes:          "Streaming: `" + channelInfo.GameName + "` \\n>>> " + userInfo.Description,
			EmbedURL:          "https://www.twitch.tv/" + twitchAccount,
			EmbedThumbnailURL: userInfo.ProfileImageURL,
		})

		w.WriteHeader(http.StatusNoContent)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("NOTIFICATION", string(bodyBuf))
		return
	case twitchAPI.MESSAGE_TYPE_VERIFICATION:
		var challengeMsg string
		if challengeMsg, err = js.Get("challenge").String(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			fmt.Println(err)
			return
		}
		// 成功 VERIFICATION
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, challengeMsg)
		fmt.Println("verification success", challengeMsg)
		return
	case twitchAPI.MESSAGE_TYPE_REVOCATION:
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("REVOCATION", string(bodyBuf))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println("unknown message type:", string(bodyBuf))
}

func TestPage(w http.ResponseWriter, r *http.Request) {
	os.Setenv("TEST", strconv.Itoa(rand.Int()))
	info := twitchAPI.UserInfo{}
	info.Validation()
	err := discordAPI.SendMessage(discordAPI.MessageOption{
		TagEveryone:       false,
		Content:           "** ﾚ(ﾟ∀ﾟ;)ﾍ 單兵注意 **`" + info.DisplayName + "`** 開直播囉 ﾍ( ﾟ∀ﾟ;)ﾉ **",
		EmbedEnable:       true,
		EmbedTitle:        info.DisplayName,
		EmbedDes:          info.Description,
		EmbedURL:          "https://www.twitch.tv/" + info.Login,
		EmbedThumbnailURL: info.ProfileImageURL,
	})
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprint(w, "Hello World")
}
