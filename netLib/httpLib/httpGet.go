package httpLib

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Get(webUrl, param string, headerMap map[string]string, transportPtr *http.Transport) (statusCode int, respBody []byte, err error) {
	if param != "" {
		if strings.Contains(webUrl, "?") {
			webUrl = fmt.Sprintf("%s&%s", webUrl, param)
		} else {
			webUrl = fmt.Sprintf("%s?%s", webUrl, param)
		}
	}

	var requestPtr *http.Request
	requestPtr, err = http.NewRequest("GET", webUrl, nil)
	if err != nil {
		return
	}

	if headerMap != nil {
		for k, v := range headerMap {
			requestPtr.Header.Add(k, v)
		}
	}

	httpClientPtr := &http.Client{}
	if transportPtr != nil {
		httpClientPtr.Transport = transportPtr
	}

	var responsePtr *http.Response
	responsePtr, err = httpClientPtr.Do(requestPtr)
	if err != nil {
		return
	}
	defer responsePtr.Body.Close()

	statusCode = responsePtr.StatusCode
	respBody, err = io.ReadAll(responsePtr.Body)

	return
}

func GetWithParamMap(webUrl string, paramMap map[string]string, header map[string]string, transportPtr *http.Transport) (statusCode int, respBody []byte, err error) {
	param := AssembleRequestString(paramMap)
	statusCode, respBody, err = Get(webUrl, param, header, transportPtr)

	return
}
