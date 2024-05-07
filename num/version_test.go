package num

import (
	"testing"
)

func equalVersion(t *testing.T, valExpect int, v1, v2 string) {
	a := NewVersion(v1)
	b := NewVersion(v2)
	valActural := a.Compare(b)
	if valActural != valExpect {
		t.Fatalf("期望: %v, 真实%v, v1 = %v, v2 = %v", valExpect, valActural, v1, v2)
	}
}
func TestVersion(t *testing.T) {
	arr := []any{
		[]any{-1, "0.0.1", "0.1"},
		[]any{-1, "0.0.1", "1"},
		[]any{-1, "0.0.0.1", "0.0.1"},

		[]any{0, "0.0.1", "0.0.1"},
		[]any{0, "0.0.1", "0.0.1.0"},
		[]any{0, "0.0.1", "0.0.1.1"},
	}
	for _, v := range arr {
		data := v.([]any)
		equalVersion(t, data[0].(int), data[1].(string), data[2].(string))
	}
}
