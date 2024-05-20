package main

import (
	srv "1linkshare/internal/server"
	"1linkshare/internal/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := utils.DeleteAllFilesInUploadDir(); err != nil {
		log.Println("Error cleaning up upload directory:", err)
	}
}

func main() {

	log.Println("Application started")
	port := srv.GetPortOrDefault(8080)
	server := srv.NewServer(srv.WithPort(port))
	server.Start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
}
