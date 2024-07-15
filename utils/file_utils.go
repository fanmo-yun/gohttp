package utils

import (
	"os"

	"go.uber.org/zap"
)

func VerifyPath(path string) (bool, error) {
	f, openErr := os.OpenFile(path, os.O_RDONLY, 0444)
	if openErr != nil {
		zap.L().Error("Path verification failed", zap.String("path", path), zap.Error(openErr))
		return false, openErr
	}
	defer f.Close()

	fileInfo, statErr := f.Stat()
	if statErr != nil {
		zap.L().Error("Path verification failed", zap.String("path", path), zap.Error(statErr))
		return false, statErr
	}

	return fileInfo.IsDir(), nil
}
