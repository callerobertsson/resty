package utils

import (
	"errors"
	"os"
)

func FileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}
