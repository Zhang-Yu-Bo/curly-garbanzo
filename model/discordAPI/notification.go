package discordAPI

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type MessageOption struct {
	TagEveryone       bool
	Content           string
	EmbedEnable       bool
	EmbedTitle        string
	EmbedDes          string
	EmbedURL          string
	EmbedThumbnailURL string
}

const templateMessage = `
{
	"content": "{{ if .TagEveryone }}@everyone {{ end }}{{ .Content }}"
	{{ if .EmbedEnable }}
	,"embeds": 
	[
		{
			"title": "{{ .EmbedTitle }}",
			"description": "{{ .EmbedDes }}",
			"url": "{{ .EmbedURL }}",
			"thumbnail": {
                "url": "{{ .EmbedThumbnailURL }}"
            }
		}
	]
	{{ end }}
	{{ if .TagEveryone }}
	,"allowed_mentions": {
    	"parse": ["everyone"]
  	}
	{{ end }}
}`

func SendMessage(msg MessageOption) error {
	url := os.Getenv("DISCORD_WEBHOOK")

	fmt.Println("discord webhook", url)

	var err error
	var bufMsg bytes.Buffer
	templateMsg := template.Must(template.New("dcMsg").Parse(templateMessage))
	if err = templateMsg.Execute(&bufMsg, msg); err != nil {
		return err
	}
	payload := strings.NewReader(bufMsg.String())

	var req *http.Request
	client := &http.Client{}

	if req, err = http.NewRequest("POST", url, payload); err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		resultMsg, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(resultMsg))
	}

	return nil
}
