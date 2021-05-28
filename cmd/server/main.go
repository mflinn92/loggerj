package main

import (
	"log"

	"github.com/mflinn92/loggerj/internal/server"
)

func main() {
	server := server.NewHTTPServer(":8000", server.NewLog())
	log.Fatal(server.ListenAndServe())
}
