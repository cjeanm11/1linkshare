package utils

import (
	"log"
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


func DeleteAllFilesInUploadDir() error {
	uploadDir := "./uploads"
	files, _ := os.ReadDir(uploadDir)
	for _, file := range files {
		filePath := filepath.Join(uploadDir, file.Name())
		if err := os.Remove(filePath); err != nil {
			log.Println("Error deleting file:", err)
		}
	}
	return nil
}