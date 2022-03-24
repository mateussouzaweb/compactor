package os

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

// RandomString generates a random string from give size
func RandomString(n int) string {

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// TemporaryFile return a temporary file path
func TemporaryFile(file string) string {
	fileName := Name(file) + "-" + RandomString(10) + Extension(file)
	return filepath.Join(os.TempDir(), fileName)
}
