package sass

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

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
	file := system.TemporaryFile("sass-transpiler.js")
	defer func() {
		errors.Join(err, system.Delete(file))
	}()

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

	return nil
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
func (service *TranspilerService) Execute(config *SassConfig, file *processor.File) error {

	data := struct {
		Config   *SassConfig `json:"config"`
		Source   string      `json:"source"`
		Filename string      `json:"filename"`
		Relative string      `json:"relative"`
	}{
		Config:   config,
		Source:   file.Content,
		Filename: file.File,
		Relative: system.Relative(system.Dir(file.Destination), file.Path),
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

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result := struct {
		Output    string `json:"output"`
		SourceMap string `json:"sourceMap"`
	}{}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return err
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

	return nil
}
