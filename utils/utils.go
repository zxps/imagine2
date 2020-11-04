package utils

import (
	"math/rand"
	"os"
)

// IsFileExists - check if file exists
func IsFileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// IsDirExists ...
func IsDirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// RandInt - generate random integer in range from "min" to "max"
func RandInt(min int, max int) int {
	result := rand.Intn(max-min) + min

	return result
}
