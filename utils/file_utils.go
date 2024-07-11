package utils

import (
	"os"
)

func VerifyPath(path string) (bool, error) {
	f, openErr := os.OpenFile(path, os.O_RDONLY, 0444)
	if openErr != nil {
		return false, openErr
	}
	defer f.Close()

	fileInfo, statErr := f.Stat()
	if statErr != nil {
		return false, statErr
	}

	return fileInfo.IsDir(), nil
}
