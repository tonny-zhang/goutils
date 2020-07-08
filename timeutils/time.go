package timeutils

import "time"

// ParseWithLocation 解析为当前时区时间
func ParseWithLocation(timeFormat, timeStr string) (t time.Time, err error) {
	locationName := "Asia/Shanghai"
	// https://studygolang.com/articles/14933?fr=sidebar
	l, err := time.LoadLocation(locationName)
	if err != nil {
		return
	}
	t, err = time.ParseInLocation(timeFormat, timeStr, l)
	return
}