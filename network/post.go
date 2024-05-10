package network

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// PostJSON send post with json
func PostJSON(urlPost string, params any, headers ...Header) (body []byte, e error) {
	headers = append(headers, Header{
		Key:   "Content-Type",
		Value: "application/json;charset=UTF-8",
	})
	return Post(urlPost, params, headers...)
}

// Post 发送post请求
func Post(urlPost string, params any, headers ...Header) (body []byte, e error) {
	bytesData, e := json.Marshal(params)
	if e != nil {
		return
	}
	// fmt.Println(string(bytesData))
	reader := bytes.NewReader(bytesData)
	request, e := http.NewRequest("POST", urlPost, reader)
	if e != nil {
		return
	}

	var useProxy = false
	for _, header := range headers {
		if header.Key == "use_proxy" {
			useProxy = true
			continue
		}
		request.Header.Set(header.Key, header.Value)
	}

	client := http.Client{}

	if useProxy {
		client.Transport = GetTransport()
	}

	resp, e := client.Do(request)
	if e != nil {
		return
	}
	body, e = io.ReadAll(resp.Body)

	return
}

// GetJSONWithPostJSON post请求得到json结构
func GetJSONWithPostJSON(urlPost string, target any, params any, headers ...Header) (e error) {
	b, e := PostJSON(urlPost, params, headers...)
	if e != nil {
		return
	}

	// fmt.Println("响应结果：", string(b))
	e = json.Unmarshal(b, target)
	return
}
