package util

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Ternary is a generic function that returns one of two values based on a condition.
func Ternary[T any](cond bool, trueVal T, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

func DefaultIfZero(num int, defaultNum int) int {
	if num == 0 {
		return defaultNum
	}
	return num
}

func StringVal(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func JoinNumbers(nums []int64) string {
	strs := make([]string, len(nums))

	for i, num := range nums {
		strs[i] = strconv.FormatInt(num, 10)
	}

	return strings.Join(strs, ", ")
}

// PrettyPrint formats JSON data with indentation.
func PrettyPrint(data interface{}) string {
	if data == nil {
		return ""
	}

	d, _ := json.MarshalIndent(data, "", "  ")
	return string(d)
}
