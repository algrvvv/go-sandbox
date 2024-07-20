package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func main() {
	codePath := os.Args[1]
	code, err := os.ReadFile(codePath)
	if err != nil {
		log.Printf("failed to read code: %v", err)
		return
	}

	cmd := exec.Command("go", "run", "-")
	cmd.Stdin = bytes.NewReader(code)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		log.Printf("failed to run code: %v", err)
	}
}
