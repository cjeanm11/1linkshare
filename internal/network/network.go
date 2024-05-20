package network

import (
	"bufio"
	"fmt"
	"os/exec"
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
