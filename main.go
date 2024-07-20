package main

import (
	"github.com/algrvvv/go-sandbox/src"
	"log"
	"net/http"
)

func main() {
	serv := src.GetNewServer()
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", serv))
}
