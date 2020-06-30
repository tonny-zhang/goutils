package cache

import (
	"fmt"
	"testing"
)

func TestGetCache(t *testing.T) {
	b, e := GetCache("/Users/mac/source/projects/bitgene_utils/file/file.go", 0)
	fmt.Println(b, e)
}
