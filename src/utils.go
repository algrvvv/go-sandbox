package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
		return ExecutionResult{code, "", err.Error()}
	}
	defer os.Remove(codeFilePath)

	//cmd := exec.Command("docker", "run", "--rm", "-v", "/tmp/go-sandbox:/app",
	//	"-v", "/var/run/docker.sock:/var/run/docker.sock",
	//	"--security-opt", "seccomp=seccomp.json",
	//	"go-runner", "go", "run", "/app/"+fileName)

	cmd := exec.Command("docker", "run", "--rm", "-v", "/tmp/go-sandbox:/app",
		"go-runner", "go", "run", "/app/"+fileName)

	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Println("start execution err: ", err)
		return ExecutionResult{code, out.String(), err.Error()}
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		return ExecutionResult{code, out.String(), "Timeout: code execution took too long"}
	case err := <-done:
		if err != nil {
			log.Println("execution err: ", err.Error())
			errMsg := fmt.Sprintf("%s\n%s", stderr.String(), err.Error())
			return ExecutionResult{code, out.String(), errMsg}
		}
	}

	return ExecutionResult{code, out.String(), stderr.String()}
}
