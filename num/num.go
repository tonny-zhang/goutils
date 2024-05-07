package num

import (
	"strconv"
	"strings"
)

// Join join int array
func Join[T byte | int | int16 | int32 | int64](slice []T, dep string) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(int(v))
	}
	return strings.Join(strSlice, dep)
}

// Split strings split to int slice
func Split[T byte | int | int16 | int32 | int64](str, sep string) (res []T) {
	strSlice := strings.Split(str, sep)
	l := len(strSlice)
	res = make([]T, l)

	for i := 0; i < l; i++ {
		if v, e := strconv.Atoi(strSlice[i]); e == nil {
			res[i] = T(v)
		}
	}

	return res
}
