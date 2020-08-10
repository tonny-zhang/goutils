package network

import (
	"bytes"
	"fmt"
	"goutils/encode"
	"goutils/fileutils"

	"io"
	"io/ioutil"
	"net/http"
	libURL "net/url"
	"os"
	"path"
)

// Download 下载文件
func Download(url, outdir string, extDefault string) (outfile string, err error) {
	var ext string
	if extDefault == "" {
		u, e := libURL.Parse(url)
		if e == nil {
			ext = path.Ext(u.Path)
		}
		if ext == "" {
			ext = ".png"
		}
	} else {
		ext = extDefault
	}

	// 目录不存在时创建
	if !fileutils.IsFileExists(outdir) {
		os.MkdirAll(outdir, os.ModePerm)
	}

	resp, err := http.Get(url)

	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s 访问错误", url)
	} else {
		outfile = path.Join(outdir, encode.MD5(url)+ext)
		body, _ := ioutil.ReadAll(resp.Body)
		out, _ := os.Create(outfile)
		io.Copy(out, bytes.NewReader(body))
	}

	return
}
