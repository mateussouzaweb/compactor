package compactor

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecCommand(cmd string, args ...string) (string, error) {

	output, err := exec.Command(cmd, args...).Output()

	if err != nil {
		_args := strings.Join(args, " ")
		return "", fmt.Errorf("Command error: %s %s\n%v\n%s", cmd, _args, err, string(output))
	}

	return string(output), nil
}
