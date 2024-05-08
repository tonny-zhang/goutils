package fileutils

import (
	"fmt"
	"testing"
)

func equalParseRuntimeProjectBaseDir(t *testing.T, desc string, stack string, mainModuleName string, expect string) {
	res := parseRuntimeProjectBaseDir(stack, mainModuleName)
	if res != expect {
		t.Errorf("%s: %s", desc, fmt.Sprintf("expect: %s, but got: %s", expect, res))
	} else {
		// t.Logf("%s: %s", desc, fmt.Sprintf("expect: %s, got: %s", expect, res))
	}
}
func TestParseRuntimeProjectBaseDir(t *testing.T) {
	equalParseRuntimeProjectBaseDir(t, "main方法里直接调用",
		`
main.main()
	/project/abc/main.go:30 +0x6c
`, "abc", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "main协程方法里直接调用",
		`
main.main.func1()
	/project/abc/main.go:30 +0x6c
created by main.main in goroutine 1
	/project/abc/main.go:32 +0x1a
`, "abc", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "子目录方法里直接调用",
		`
abc/lib.Run()
	/project/abc/lib/test.go:30 +0x6c
`, "abc", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "子目录协程方法里直接调用",
		`
abc/lib.Run.func1()
	/project/abc/lib/test.go:30 +0x6c
created by abc/lib.Run in goroutine 1
	/project/abc/lib/test.go:32 +0x1a
`, "abc", "/project/abc")

	// ----------目录名和模板名不一致

	// fmt.Println(" ----------目录名和模板名不一致")
	equalParseRuntimeProjectBaseDir(t, "目录名和模板名不一致 main方法里直接调用",
		`
main.main()
	/project/abc/main.go:30 +0x6c
`, "abc1", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "目录名和模板名不一致 main协程方法里直接调用",
		`
main.main.func1()
	/project/abc/main.go:30 +0x6c
created by main.main in goroutine 1
	/project/abc/main.go:32 +0x1a
`, "abc1", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "目录名和模板名不一致 子目录方法里直接调用",
		`
abc1/lib.Run()
	/project/abc/lib/test.go:30 +0x6c
`, "abc1", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "目录名和模板名不一致 子目录协程方法里直接调用",
		`
abc1/lib.Run.func1()
	/project/abc/lib/test.go:30 +0x6c
created by abc1/lib.Run in goroutine 1
	/project/abc/lib/test.go:32 +0x1a
`, "abc1", "/project/abc")

}
