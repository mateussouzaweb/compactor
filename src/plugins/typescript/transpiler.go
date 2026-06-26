package typescript

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/mateussouzaweb/compactor/src/errors"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

//go:embed *.js
var transpilerFS embed.FS

// TranspilerService struct
type TranspilerService struct {
	File    string
	Port    string
	Address string
	Cmd     *exec.Cmd
}

// Init service to handle transpilation requests
func (service *TranspilerService) Init() error {

	var err error

	// Write server script to temporary file
	file := system.TemporaryFile("typescript-transpiler.js")
	defer errors.Join(&err, func() error {
		return system.Delete(file)
	})

	serverScript, err := transpilerFS.ReadFile("transpiler.js")
	if err != nil {
		return err
	}

	err = system.Write(file, string(serverScript), 0775)
	if err != nil {
		return err
	}

	// Get temporary port
	port, err := system.TemporaryPort()
	if err != nil {
		return err
	}

	// Run server in background
	cmd := exec.Command("node", file)
	cmd.Env = append(os.Environ(), "PORT="+port)

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() error {
		err := cmd.Wait()
		if err != nil {
			panic(err)
		}
		return nil
	}()

	// Wait service become online
	address := "http://localhost:" + port
	for {
		_, err := http.Get(address)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Set service data
	service.File = file
	service.Port = port
	service.Address = address
	service.Cmd = cmd

	return err
}

// Shutdown transpilation service
func (service *TranspilerService) Shutdown() error {

	if service.Cmd.Process.Pid != 0 {
		err := service.Cmd.Process.Kill()
		if err != nil {
			return err
		}
	}

	if service.File != "" {
		err := system.Delete(service.File)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute transpilation process
func (service *TranspilerService) Execute(config *TSConfig, file *processor.File) error {

	relative := system.Relative(system.Dir(file.Destination), file.Path)
	data := struct {
		Config   *TSConfig       `json:"config"`
		File     *processor.File `json:"file"`
		Relative string          `json:"relative"`
	}{
		Config:   config,
		File:     file,
		Relative: relative,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := http.Post(
		service.Address,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	defer errors.Join(&err, func() error {
		return response.Body.Close()
	})

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result := struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		Output    string `json:"output"`
		SourceMap string `json:"sourceMap"`
	}{}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("transpiler error: %s", result.Message)
	}

	err = system.Write(file.Destination, result.Output, file.Permission)
	if err != nil {
		return err
	}

	if result.SourceMap != "" {
		err := system.Write(file.Destination+".map", result.SourceMap, file.Permission)
		if err != nil {
			return err
		}
	}

	return err
}
