package utils

import "strings"

const (
	ZeroString = ""
)

// StringTrimSpace -- trim space of string
func StringTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsStringEmpty -- check if string is empty
func IsStringEmpty(s string) bool {
	return s == ZeroString
}
