package sendMessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TgJsonData struct {
	Chat_id string `json:"chat_id"`
	Text    string `json:"text"`
}

func SendTelegramMessage(message string, chatId string, token string) error {

	const maxMessageLength = 4096
	tgUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	for len(message) > 0 {
		chunk := message
		if len(message) > maxMessageLength {
			chunk = message[:maxMessageLength]
			//중간에 짤리면 한글 깨짐. 마지막 공백을 찾아서 메세지 분리
			lastSpace := strings.LastIndex(chunk, " ")
			if lastSpace == -1 || lastSpace == 0 {
				message = message[maxMessageLength:]
			} else {

				chunk = message[:lastSpace]
				message = message[lastSpace+1:]
			}
		} else {
			chunk = message
			message = ""
		}

		tgJsonData := TgJsonData{
			Chat_id: chatId,
			Text:    chunk,
		}

		tgJsonBody, err := json.Marshal(tgJsonData)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %v", err)
		}

		resp, err := http.Post(tgUrl, "application/json", bytes.NewBuffer(tgJsonBody))
		if err != nil {
			return fmt.Errorf("http.Post failed: %v", err)
		}
		resp.Body.Close()
	}

	return nil
}

func SendLineMessage(message string, token string) error {

	lineUrl := "https://notify-api.line.me/api/notify"
	const maxMessageLength = 800 //"더보기" 안나오게 끊기
	for len(message) > 0 {
		chunk := message
		if len(message) > maxMessageLength {
			chunk = message[:maxMessageLength]
			//중간에 짤리면 한글 깨짐. 마지막 공백을 찾아서 메세지 분리
			lastSpace := strings.LastIndex(chunk, " ")
			if lastSpace == -1 || lastSpace == 0 {
				message = message[maxMessageLength:]
			} else {

				chunk = message[:lastSpace]
				message = message[lastSpace+1:]
			}
		} else {
			chunk = message
			message = ""
		}

		data := url.Values{
			"message": {chunk},
		}

		req, err := http.NewRequest("POST", lineUrl, strings.NewReader(data.Encode()))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	return nil
}
