package os

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
)

// Exec run command with given arguments
func Exec(cmd string, args ...string) (string, error) {

	result := exec.Command(cmd, args...)
	output, err := result.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("command error: %s ...\n%v\n%s", cmd, err, string(output))
	}

	return string(output), nil
}

// Checksum retrive the checksum for given content
func Checksum(content string) (string, error) {

	sum := md5.New()
	_, err := io.WriteString(sum, content)

	if err != nil {
		return "", err
	}

	inBytes := sum.Sum(nil)[:8]
	hash := hex.EncodeToString(inBytes)

	return hash, err
}
