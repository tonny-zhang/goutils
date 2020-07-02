package task

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
)

// getMD5 md5加密
func getMD5(data interface{}) string {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(data)

	h := md5.New()
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}
