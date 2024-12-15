// Package utils string and string slice funcs
package utils

// ContainsString checks if a string is member of a slice
func ContainsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
