package main

import (
	"log"

	"github.com/mflinn92/loggerj/internal/server"
)

func main() {
	server := server.NewHTTPServer(":8000")
	log.Fatal(server.ListenAndServe())
}
