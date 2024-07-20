package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func main() {
	log.Println("[Runner] start new runner")
	codePath := os.Args[1]
	log.Printf("[Runner] path: %s", codePath)
	code, err := os.ReadFile(codePath)
	if err != nil {
		log.Fatalf("failed to read code: %v", err)
	}

	cmd := exec.Command("go", "run", "-")
	cmd.Stdin = bytes.NewReader(code)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		log.Fatalf("failed to run code: %v", err)
	}
}
