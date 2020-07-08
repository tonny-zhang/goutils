package timeutils

import (
	"fmt"
	"testing"
)

func TestParseWithLocation(t *testing.T) {
	strTime := "2020-07-08 11:27:52"
	timeTo, _ := ParseWithLocation("2006-01-02 15:04:05", strTime)

	fmt.Println(timeTo)
}
