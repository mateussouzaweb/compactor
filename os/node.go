package os

import "os/exec"

// Check if node package command exists, otherwise, try to install the package
func NodeRequire(command string, pkg string) error {

	_, err := exec.LookPath(command)

	if err != nil {
		_, err = Exec("npm", "install", "-g", pkg)
	}

	return err
}
