package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func Post(url string, payload map[string]string, headers http.Header) ([]byte, error) {
	return process(http.MethodPost, url, payload, headers)
}

func Put(url string, payload map[string]string, headers http.Header) ([]byte, error) {
	return process(http.MethodPut, url, payload, headers)
}

func Get(url string, queryParams map[string]string, headers http.Header) ([]byte, error) {
	var endpoint = AppendQueryParams(url, queryParams)
	httpClient := http.Client{
		Timeout: time.Second * ConfigurationObj.HttpTimeoutInSeconds,
	}
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		ZapLoggerObj.Error("Unable to make request to endpoint " + endpoint)
		return nil, err
	}

	request.Header = headers
	response, responseErr := httpClient.Do(request)
	if responseErr != nil {
		ZapLoggerObj.Error("get request to endpoint failed " + endpoint)
		return nil, responseErr
	}

	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		ZapLoggerObj.Error("unable to convert to an object " + endpoint)
		return nil, readErr
	}
	return body, readErr
}

func process(httpMethod string, url string, payload map[string]string, header http.Header) ([]byte, error) {
	jsonPayload, _ := json.Marshal(payload)
	request, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		ZapLoggerObj.Error("[" + httpMethod + "]" + " Unable to create request to endpoint " + url +
			" with payload " + string(jsonPayload))
		return nil, err
	}
	request.Header = header
	if request.Header.Get("Content-Type") == "" {
		request.Header.Add("Content-Type", "application/json")
	}

	httpClient := http.Client{
		Timeout: time.Second * time.Duration(ConfigurationObj.HttpTimeoutInSeconds),
	}

	response, responseErr := httpClient.Do(request)
	if responseErr != nil {
		ZapLoggerObj.Error("[" + httpMethod + "]" + " request to endpoint " + url + " failed")
		return nil, responseErr
	}

	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		ZapLoggerObj.Error("[" + httpMethod + "]" + " Unable to parse response body received from " + url)
		return nil, readErr
	}
	return body, readErr
}

