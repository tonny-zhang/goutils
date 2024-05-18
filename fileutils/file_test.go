package fileutils

import (
	"fmt"
	"testing"
)

func equalParseRuntimeProjectBaseDir(t *testing.T, desc string, stack string, mainModuleName string, expect string) {
	res := parseRuntimeProjectBaseDir(stack, mainModuleName)
	if res != expect {
		t.Errorf("%s: %s", desc, fmt.Sprintf("expect: [%s], but got: [%s]", expect, res))
	} else {
		// t.Logf("%s: %s", desc, fmt.Sprintf("expect: %s, got: %s", expect, res))
	}
}
func TestParseRuntimeProjectBaseDir(t *testing.T) {
	equalParseRuntimeProjectBaseDir(t, "1目录名和模板名不一致 子目录方法里直接调用",
		`
goroutine 1 [running, locked to thread]:
runtime/debug.Stack()
		/usr/local/go/src/runtime/debug/stack.go:24 +0x5e
github.com/tonny-zhang/goutils/fileutils.initData.func1()
		/Users/tonny/source/project/goutils/fileutils/file.go:127 +0x1d1
sync.(*Once).doSlow(0xc000114580?, 0xc000107d70?)
		/usr/local/go/src/sync/once.go:74 +0xc2
sync.(*Once).Do(...)
		/usr/local/go/src/sync/once.go:65
github.com/tonny-zhang/goutils/fileutils.initData()
		/Users/tonny/source/project/goutils/fileutils/file.go:76 +0x2c
github.com/tonny-zhang/goutils/fileutils.GetCmdDir(...)
		/Users/tonny/source/project/goutils/fileutils/file.go:144
github.com/tonny-zhang/goutils/env.AutoLoad.func1()
		/Users/tonny/source/project/goutils/env/env.go:19 +0x17
sync.(*Once).doSlow(0xc000107e10?, 0x6aa127?)
		/usr/local/go/src/sync/once.go:74 +0xc2
sync.(*Once).Do(...)
		/usr/local/go/src/sync/once.go:65
github.com/tonny-zhang/goutils/env.AutoLoad()
		/Users/tonny/source/project/goutils/env/env.go:16 +0x2c
abc/lib/jwt.init.0()
		/project/abc/lib/jwt/jwt.go:27 +0x13
	`, "abc", "/project/abc")

	equalParseRuntimeProjectBaseDir(t, "2目录名和模板名不一致 子目录方法里直接调用",
		`
goroutine 1 [running]:
runtime/debug.Stack()
		/usr/local/go/src/runtime/debug/stack.go:24 +0x5e
github.com/tonny-zhang/goutils/fileutils.initData()
		/Users/tonny/source/project/goutils/fileutils/file.go:130 +0x26a
github.com/tonny-zhang/goutils/fileutils.GetRuntimeProjectBaseDir(...)
		/Users/tonny/source/project/goutils/fileutils/file.go:141
github.com/tonny-zhang/goutils/logger.Logger.Error({{0x0, 0x0}, {0xe144ea8, 0xc0000b4020}, 0x0}, {0xe025c18, 0x8}, {0x0, 0x0, 0x0})
		/Users/tonny/source/project/goutils/logger/logger.go:72 +0x90
abc/a%2eb.Run(...)
		/project/abc/a.b/test.go:7
main.main()
		/project/abc/main.go:38 +0x26f
	`, "abc", "/project/abc")

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
