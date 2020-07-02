package task

import (
	"fmt"
	"testing"
)

func TestGetMD5(t *testing.T) {
	conf := Conf{
		DirCache: "/test/abc/",
	}
	result := getMD5(conf)
	fmt.Println(result)

	fmt.Println(getMD5("test"))
}
