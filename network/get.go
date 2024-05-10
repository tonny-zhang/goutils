package network

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// GetJSON 得到json结构
func GetJSON(urlGet string, target any) (e error) {
	b, e := GetWithHeader(urlGet)
	if e != nil {
		return
	}
	e = json.Unmarshal(b, target)
	return
}

// GetJSONWithHeader 得到json结构
func GetJSONWithHeader(urlGet string, target any, headers ...Header) (e error) {
	b, e := GetWithHeader(urlGet, headers...)
	if e != nil {
		return
	}
	e = json.Unmarshal(b, target)
	return
}

// GetWithHeader Get请求
func GetWithHeader(urlGet string, headers ...Header) (body []byte, err error) {
	URL, err := url.Parse(urlGet)
	if err != nil {
		return
	}

	urlPath := URL.String()

	request, err := http.NewRequest("GET", urlPath, nil)

	for _, header := range headers {
		request.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	return content, nil
}
