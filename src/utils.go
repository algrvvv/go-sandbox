package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func getHashedFileName(fileName string) string {
	h := sha256.New()
	h.Write([]byte(fileName))
	return hex.EncodeToString(h.Sum(nil))
}

func executeUserCode(code string) ExecutionResult {
	fileName := getHashedFileName(time.Now().Format("20060102150405")) + ".go"
	codeFilePath := filepath.Join("/tmp/go-sandbox", fileName)
	if err := os.WriteFile(codeFilePath, []byte(code), 0644); err != nil {
		log.Println("write file err: ", err)
		return ExecutionResult{"", err.Error()}
	}

	cmd := exec.Command("docker", "run", "--rm", "-v", "/tmp/go-sandbox:/app",
		//"--security-opt", "seccomp=seccomp.json",
		"go-runner", "go", "run", "/app/"+fileName)

	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Println("start execution err: ", err)
		return ExecutionResult{"", err.Error()}
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		return ExecutionResult{"", "Timeout: code execution took too long"}
	case err := <-done:
		if err != nil {
			log.Println("execution err: ", err)
			return ExecutionResult{out.String(), err.Error()}
		}
	}

	return ExecutionResult{out.String(), stderr.String()}
}
