package os

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Exec run command with given arguments
func Exec(cmd string, args ...string) (string, error) {

	output, err := exec.Command(cmd, args...).Output()

	if err != nil {
		_args := strings.Join(args, " ")
		return "", fmt.Errorf("command error: %s %s\n%v\n%s", cmd, _args, err, string(output))
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
