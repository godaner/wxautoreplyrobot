package go_util

import "strings"

func TrimRightSpace(s string) string {
	return strings.TrimRight(string(s), "\r\n\t ")
}