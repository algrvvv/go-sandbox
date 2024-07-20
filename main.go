package main

import (
	"log"
	"net/http"

	"github.com/algrvvv/go-sandbox/src"
)

func main() {
	serv := src.GetNewServer()
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", serv))
}
