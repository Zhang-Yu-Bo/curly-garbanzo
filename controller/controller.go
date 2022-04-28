package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zhang-Yu-Bo/curly-garbanzo/model/discordAPI"
	"github.com/Zhang-Yu-Bo/curly-garbanzo/model/twitchAPI"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func ShowUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var err error
	var userInfo twitchAPI.UserInfo

	userInfo, err = twitchAPI.GetUserInfoByName(vars["username"], twitchAPI.GetAppAccessToken())
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
		fmt.Println("hmac =", hmac)
		fmt.Println("1234 =", r.Header.Get(twitchAPI.TWITCH_MESSAGE_SIGNATURE))
		fmt.Println("validation failed, hmac is not equal to TWITCH_MESSAGE_SIGNATURE")
		return
	}

	switch r.Header.Get(twitchAPI.MESSAGE_TYPE) {
	case twitchAPI.MESSAGE_TYPE_NOTIFICATION:

		discordAPI.SendMessage(discordAPI.MessageOption{
			TagEveryone:       true,
			Content:           "test",
			EmbedEnable:       true,
			EmbedTitle:        "this is title",
			EmbedDes:          "this is description",
			EmbedURL:          "https://www.twitch.tv/beryl_lulu",
			EmbedThumbnailURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/6709d95f-3cbc-4f2b-920f-3b408be0dc96-profile_image-70x70.png",
		})

		w.WriteHeader(http.StatusNoContent)
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
	err := discordAPI.SendMessage(discordAPI.MessageOption{
		TagEveryone:       true,
		Content:           "test",
		EmbedEnable:       true,
		EmbedTitle:        "this is title",
		EmbedDes:          "this is description",
		EmbedURL:          "https://www.twitch.tv/beryl_lulu",
		EmbedThumbnailURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/6709d95f-3cbc-4f2b-920f-3b408be0dc96-profile_image-70x70.png",
	})
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprint(w, "Hello World")
}
