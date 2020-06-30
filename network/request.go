package network

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// Get Get请求
func Get(urlGet string, params map[string]string) (body []byte, err error) {
	paramVal := url.Values{}
	URL, e1 := url.Parse(urlGet)
	if e1 != nil {
		return nil, e1
	}

	for k, v := range params {
		paramVal.Set(k, v)
	}
	URL.RawQuery = paramVal.Encode()

	urlPath := URL.String()

	resp, e2 := http.Get(urlPath)
	if e2 != nil {
		return nil, e2
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return content, nil
}

// GetWithHeader Get请求
func GetWithHeader(urlGet string) (body []byte, err error) {
	URL, err := url.Parse(urlGet)
	if err != nil {
		return
	}

	urlPath := URL.String()

	request, err := http.NewRequest("GET", urlPath, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return content, nil
}

// Post Post请求
func Post(urlPost string, params map[string]string) (body []byte, err error) {
	paramVal := url.Values{}
	URL, e1 := url.Parse(urlPost)
	if e1 != nil {
		return nil, e1
	}

	for k, v := range params {
		paramVal.Set(k, v)
	}

	urlPath := URL.String()

	resp, e2 := http.PostForm(urlPath, paramVal)

	if e2 != nil {
		return nil, e2
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("出现错误, 响应码[%d]", resp.StatusCode)
	} else {
		content, _ := ioutil.ReadAll(resp.Body)

		body = content
	}

	return
}

// PostUpload 上传文件
func PostUpload(urlPost string, params map[string]string, uploadParamName, uploadFilePath string) (body []byte, err error) {
	postData := &bytes.Buffer{}
	writer := multipart.NewWriter(postData)

	if uploadFilePath != "" {
		file, err := os.Open(uploadFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(uploadParamName, uploadFilePath)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
		// fmt.Printf("%s = %s\n", key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", urlPost, postData)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("error to request to the server:%s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	return content, nil
}
