package network

import (
	"bytes"
	"goutils/encode"
	"goutils/fileutils"

	"io"
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

	body, err := GetWithHeader(url)
	if err == nil {
		// 目录不存在时创建
		if !fileutils.IsFileExists(outdir) {
			os.MkdirAll(outdir, os.ModePerm)
		}
		outfile = path.Join(outdir, encode.MD5(url)+ext)
		out, _ := os.Create(outfile)
		io.Copy(out, bytes.NewReader(body))
	}

	// resp, err := http.Get(url)

	// if err != nil {
	// 	return
	// }
	// if resp.StatusCode != 200 {
	// 	err = fmt.Errorf("%s 访问错误", url)
	// } else {
	// 	outfile = path.Join(outdir, encode.MD5(url)+ext)
	// 	body, _ := ioutil.ReadAll(resp.Body)
	// 	out, _ := os.Create(outfile)
	// 	io.Copy(out, bytes.NewReader(body))
	// }

	return
}
