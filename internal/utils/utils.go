package utils

import (
	"os"
	"path/filepath"
)

func CreateFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Create(filePath)
}