package converter

import (
	"strconv"
)

// StringToUint64
// if an error occurred, this will return 0
func StringToUint64(stringNumber string) uint64 {
	value, err := strconv.ParseInt(stringNumber, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(value)
}
func StringToInt(stringNumber string) int {
	value, err := strconv.ParseInt(stringNumber, 10, 64)
	if err != nil {
		return 0
	}
	return int(value)
}
