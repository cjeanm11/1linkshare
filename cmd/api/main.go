package main

import (
	"log"
	srv "1linkshare/internal/server"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	log.Println("Application started")
	port := srv.GetPortOrDefault(8080)
	server := srv.NewServer(
		srv.WithPort(port),
		srv.WithTSL(false),
		srv.WithGRPC(false),
	)
	server.Start()
}
