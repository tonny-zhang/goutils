package num

import (
	"reflect"
	"testing"
)

func equalJoin[T byte | int | int16 | int32 | int64](t *testing.T, arr []T, valExpect string) {
	valActural := Join(arr, "-")
	if valActural != valExpect {
		t.Fatalf("%v %v 期望: %v, 真实%v", reflect.TypeOf(arr), reflect.ValueOf(arr), valExpect, valActural)
	}
}

func TestJoin(t *testing.T) {
	equalJoin(t, []byte{1, 2, 3, 4, 5}, "1-2-3-4-5")
	equalJoin(t, []int{1, 2, 3, 4, 5}, "1-2-3-4-5")
	equalJoin(t, []int16{1, 2, 3, 4, 5}, "1-2-3-4-5")
	equalJoin(t, []int32{1, 2, 3, 4, 5}, "1-2-3-4-5")
	equalJoin(t, []int64{1, 2, 3, 4, 5}, "1-2-3-4-5")
}
func equalSplit[T byte | int | int16 | int32 | int64](t *testing.T, str string, valExpect []T) {
	valActural := Split[T](str, "-")

	tActural := reflect.TypeOf(valActural)
	tExpect := reflect.TypeOf(valExpect)

	var flag = tActural == tExpect
	if flag {
		flag = len(valActural) == len(valExpect)
		if flag {
			for i := 0; i < len(valActural); i++ {
				if valActural[i] != valExpect[i] {
					flag = false
					break
				}
			}
		}
	}
	if !flag {
		t.Fatalf("%v 期望: %v %v, 真实 %v %v", str, tExpect, valExpect, tActural, valActural)
	}
}

func TestSplit(t *testing.T) {
	equalSplit(t, "1-2-3", []byte{1, 2, 3})
	equalSplit(t, "1-2-3", []int{1, 2, 3})
	equalSplit(t, "1-2-3", []int16{1, 2, 3})
	equalSplit(t, "1-2-3", []int32{1, 2, 3})
	equalSplit(t, "1-2-3", []int64{1, 2, 3})

	// t.Fail()
}
