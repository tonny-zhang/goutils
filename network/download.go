package network

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/tonny-zhang/goutils/fileutils"
)

// DownloadWithSave 下载文件
func DownloadWithSave(url string, savePath string, useProxy bool) (e error) {
	request, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return
	}
	client := &http.Client{}
	if useProxy {
		client.Transport = GetTransport()
	}
	req, e := client.Do(request)
	if e != nil {
		return
	}

	defer req.Body.Close()

	fileutils.Mkdirp(path.Dir(savePath))

	f, e := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if e != nil {
		return
	}
	defer f.Close()
	_, e = io.Copy(f, req.Body)

	return
}
