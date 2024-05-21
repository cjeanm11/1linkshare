package network

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func RunSSHCommand(urlCh chan<- string, args []string) {
	cmd := exec.Command("ssh", args...)
	
	cmd.Stdin = nil

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting SSH command:", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, " ")
			lastword := parts[len(parts)-1]
			if isValidURL(lastword) {
				urlCh <- lastword
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for SSH command:", err)
		return
	}
}

func isValidURL(url string) bool {
	regex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`)
	return regex.MatchString(url)
}


func UploadFile(url, filePath string) error {
    cleanedPath := filepath.Clean(filePath)
	file, err := os.Open(cleanedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadFile", filepath.Base(cleanedPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}