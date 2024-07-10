package utils

import (
	"os"
	"path/filepath"
)

func VerifyFileExistence(path string) error {
	fp := filepath.Join("html", path)
	f, openErr := os.OpenFile(fp, os.O_RDWR, 0666)
	if openErr != nil {
		return openErr
	}
	defer f.Close()
	return nil
}
