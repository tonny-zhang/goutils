package str

import (
	"math/rand"
	"time"
)

var letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var letterBytesLen = len(letterBytes)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SetLetter 设置生成随机字符串的字符集
func SetLetter(letter string) {
	letterBytes = letter
	letterBytesLen = len(letterBytes)
}

// GetNoce 得到指定位数随机字符串
func GetNoce(num int) string {
	b := make([]byte, num)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(letterBytesLen)]
	}
	return string(b)
}
