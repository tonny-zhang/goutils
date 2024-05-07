package num

import (
	"strconv"
	"strings"
)

// Version 版本处理
type Version struct {
	v1 int
	v2 int
	v3 int

	score int
}

// IsEmpty 是否设置值
func (v Version) IsEmpty() bool {
	return v.score == 0
}

// Compare 比较大小
func (v Version) Compare(b Version) int {
	if v.v1 > b.v1 {
		return 1
	} else if v.v1 < b.v1 {
		return -1
	} else {
		if v.v2 > b.v2 {
			return 1
		} else if v.v2 < b.v2 {
			return -1
		} else {
			if v.v3 > b.v3 {
				return 1
			} else if v.v3 < b.v3 {
				return -1
			} else {
				return 0
			}
		}
	}
}

// NewVersion 得到一个版本描述
func NewVersion(version string) Version {
	ins := Version{}
	vArr := strings.Split(version, ".")
	lenV := len(vArr)
	if lenV > 0 {
		if v, e := strconv.Atoi(vArr[0]); e == nil {
			ins.v1 = v
		}
		if lenV > 1 {
			if v, e := strconv.Atoi(vArr[1]); e == nil {
				ins.v2 = v
			}

			if lenV > 2 {
				if v, e := strconv.Atoi(vArr[2]); e == nil {
					ins.v3 = v
				}
			}
		}
		ins.score = ((32 << ins.v1) + (16 << ins.v2) + ins.v3)
	}

	return ins
}
