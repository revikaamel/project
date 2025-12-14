package util

import "strings"

// cek apakah string kosong
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
