package utils

func Abs(value int64) uint64 {
	if value < 0 {
		return uint64(-value)
	} else {
		return uint64(value)
	}
}
