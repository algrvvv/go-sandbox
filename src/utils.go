package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/algrvvv/go-sandbox/src/logger"
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
		logger.Error("failed to write code to /tmp/go-sandbox: "+err.Error(), err)
		return ExecutionResult{code, "", err.Error()}
	}
	defer os.Remove(codeFilePath)

	// TODO "--security-opt", "seccomp=seccomp.json",
	cmd := exec.Command("docker", "run", "--rm",
		"-v", "code-files:/app",
		"go-runner", "/app/"+fileName)

	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		logger.Error("failed to start docker run: "+err.Error(), err)
		return ExecutionResult{code, out.String(), err.Error()}
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(10 * time.Second):
		logger.Warn("docker run timeout")
		_ = cmd.Process.Kill()
		return ExecutionResult{code, out.String(), "Timeout: code execution took too long"}
	case err := <-done:
		if err != nil {
			logger.Error("docker run failed: "+err.Error(), err)
			errMsg := fmt.Sprintf("%s\n%s", stderr.String(), err.Error())
			return ExecutionResult{code, out.String(), errMsg}
		}
	}

	return ExecutionResult{code, out.String(), stderr.String()}
}
