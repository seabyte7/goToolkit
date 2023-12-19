package httpLib

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func Post(webUrl string, dataBytes []byte, headerMap map[string]string, transport *http.Transport) (statusCode int, respBody []byte, err error) {
	var requestPtr *http.Request
	requestPtr, err = http.NewRequest("POST", webUrl, bytes.NewReader(dataBytes))
	if err != nil {
		return
	}

	if headerMap != nil {
		for k, v := range headerMap {
			requestPtr.Header.Add(k, v)
		}
	}

	httpClientPtr := &http.Client{}
	if transport != nil {
		httpClientPtr.Transport = transport
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

func PostWithParamMap(webUrl string, paramMap map[string]string, headerMap map[string]string, transportPtr *http.Transport) (statusCode int, respBody []byte, err error) {
	postValues := url.Values{}
	if paramMap != nil {
		for key, value := range paramMap {
			postValues.Set(key, value)
		}
	}
	statusCode, respBody, err = Post(webUrl, []byte(postValues.Encode()), headerMap, transportPtr)

	return
}
