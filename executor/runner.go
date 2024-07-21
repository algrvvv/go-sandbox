package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	codePath := os.Args[len(os.Args)-1]

	cmd := exec.Command("go", "run", codePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("failed to run code: %v", err)
	}
}
