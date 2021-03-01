package network

import (
	"encoding/json"
	"net/http"
)

// TypeJSON 响应的json结构
type TypeJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// JSON http响应json
func JSON(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// if data == nil {
	// 	data = make(map[string]interface{})
	// }
	d := TypeJSON{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	b, _ := json.Marshal(d)
	w.Write(b)
}

// CheckPost 检查post方法
func CheckPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		w.WriteHeader(404)
		w.Write([]byte("404 not found"))
		return false
	}
	// 解析请求参数
	if t := r.Header["Content-Type"]; len(t) > 0 && t[0] == "application/x-www-form-urlencoded" {
		r.ParseForm()
	} else {
		r.ParseMultipartForm(32 << 20)
	}
	return true
}
