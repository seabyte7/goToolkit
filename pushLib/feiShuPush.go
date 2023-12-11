package pushLib

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"goToolkit/logLib"
	"goToolkit/netLib/httpLib"
	"net/http"
)

/*
FeiShu document: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
*/

var (
	defaultFeiShuWebHookUrl = "" // "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx"
)

func SetDefaultFeiShuWebHookUrl(url string) {
	defaultFeiShuWebHookUrl = url
}

func IsvalidFeiShuWebHookUrl(url string) bool {
	return len(url) > 0
}

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

// push text message to feiShu default webHookUrl
func PushTextMessageToDefault(message string) error {
	return PushTextMessage(defaultFeiShuWebHookUrl, message)
}

func PushTextMessage(webHookUrl, message string) error {
	if len(webHookUrl) == 0 {
		webHookUrl = defaultFeiShuWebHookUrl
	}
	if IsvalidFeiShuWebHookUrl(webHookUrl) {
		logLib.Zap().Error("PushTextMessage invalid webHookUrl",
			zap.String("webHookUrl", webHookUrl),
			zap.String("message", message))
		return errors.New("invalid webHookUrl")
	}

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
