package mod1

import (
	"fmt"

	"github.com/tonny-zhang/goutils/mod2"
)

// MethodMod1 mod1
func MethodMod1() {
	fmt.Println("MethodMod1")

	mod2.MethodMod2()
}
