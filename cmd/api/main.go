package main

import (
	"log"
	"os"
	//	"os"
	srv "server-template/internal/server"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()
}

func main() {
	log.Println("Application started")
	port := srv.GetPortOrDefault(8080)
	server := srv.NewServer(srv.WithPort(port), srv.WithTSL(false))
	server.Start()
}
