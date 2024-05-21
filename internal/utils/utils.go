package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func isValidFilePath(filePath string) bool {
    if filePath == "" {
        return false
    }
    if _, err := filepath.Abs(filePath); err != nil {
        return false
    }
    return true
}


func CreateFile(filePath string) (*os.File, error) {
    cleanedPath := filepath.Clean(filePath)
    if !isValidFilePath(cleanedPath) {
        return nil, errors.New("invalid file path")
    }
    dir := filepath.Dir(cleanedPath)

    if err := os.MkdirAll(dir, 0750); err != nil {
        return nil, err
    }
    f, err := os.Create(cleanedPath)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func DeleteAllFilesInUploadDir() error {
	uploadDir := "./uploads"
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(uploadDir, file.Name())
		if err := os.Remove(filePath); err != nil {
			log.Print("Error deleting file:", err)
		}
	}
	return nil
}
