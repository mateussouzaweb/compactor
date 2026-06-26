package main

import (
	"os"
	"os/exec"

	"github.com/mateussouzaweb/compactor/src/cli"
)

func main() {

	// Exit with proper code
	exitCode := 0
	defer os.Exit(exitCode)

	// Create the build script
	// Will build binaries for both architectures
	script := `mkdir -p bin/; \
	export GOOS=linux; export GOARCH=amd64; \
	go build go -buildvcs=false -o bin/compactor-amd64 cmd/compactor; \
	export GOOS=linux; export GOARCH=arm64; \
	go build go -buildvcs=false -o bin/compactor-arm64 cmd/compactor`

	cmd := exec.Command("bash", "-c", script)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and check for errors
	err := cmd.Run()
	if err != nil {
		cli.Printf(cli.Fatal, "[ERROR] Build error: %s\n", err.Error())
		exitCode = 1
	} else {
		cli.Printf(cli.Success, "[SUCCESS] Build completed successfully!\n")
	}

}
