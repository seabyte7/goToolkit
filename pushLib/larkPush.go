package pushLib

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"goToolkit/netLib/httpLib"
	"net/http"
)

/*
FeiShu document: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
*/

type FeiShuTextMessage struct {
	MessageType string                   `json:"msg_type"`
	Content     FeiShuTextMessageContent `json:"content"`
}

type FeiShuTextMessageContent struct {
	Text string `json:"text"`
}

type FeiShuResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func PushTextMessage(webHookUrl, message string) error {
	msg := FeiShuTextMessage{
		MessageType: "text",
		Content: FeiShuTextMessageContent{
			Text: message,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := httpLib.GetJsonHeader()
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	statusCode, result, err := httpLib.Post(webHookUrl, data, header, transport)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", statusCode)
	}

	response := &FeiShuResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return err
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}

func PushTextMessageToAll(webHookUrl, message string) error {
	message = fmt.Sprintf("<at user_id=\"all\">everyone</at>\n%s", message)
	return PushTextMessage(webHookUrl, message)
}
