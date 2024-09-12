package main

import (
	"log"

	"github.com/algrvvv/go-sandbox/src"
	"github.com/algrvvv/go-sandbox/src/logger"
)

func main() {
	if err := logger.NewLogger("sandbox.log"); err != nil {
		log.Fatal(err)
	}

	serv := src.GetNewServer()
	logger.Info("Starting server on port 8080...")
	if err := serv.ListenAndServe(); err != nil {
		logger.Error("Server error: "+err.Error(), err)
	}
}
