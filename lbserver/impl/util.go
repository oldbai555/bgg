package impl

import "strconv"

func toUint64(v string) uint64 {
	newV, _ := strconv.ParseUint(v, 10, 64)
	return newV
}
