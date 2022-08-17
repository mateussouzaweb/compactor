package typescript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

type TranspilerService struct {
	File    string
	Port    string
	Address string
	Cmd     *exec.Cmd
}

// Init service to handle transpilation requests
func (service *TranspilerService) Init() error {

	port, err := os.TemporaryPort()

	if err != nil {
		return err
	}

	code := fmt.Sprintf(`
	const { execSync } = require("child_process")
	const root = execSync("npm root -g").toString().trim()
	const ts = require(root + "/typescript")
	const http = require("http")
	const url = require("url")

	const httpServer = http.createServer(async (request, response) => {

		try {
			
			const buffers = []
			for await (const chunk of request) {
				buffers.push(chunk)
			}

			const data = Buffer.concat(buffers).toString()
			const body = JSON.parse(data)

			const config = body.config
				config.fileName = body.relative

			const source = body.source
			const result = ts.transpileModule(source, config)

			const output = result.outputText ? result.outputText : ""
			const sourceMap = result.sourceMapText ? result.sourceMapText.replace(
				'"sources":["' + body.filename + '"]',
				'"sources":["' + body.relative + '"]'
			) : ""

			response.writeHead(200, { "Content-Type": "application/json" })
			response.write(JSON.stringify({
				output: output,
				sourceMap: sourceMap
			}))

		} catch (error){
			
			response.writeHead(400, { "Content-Type": "application/json" })
			response.write(JSON.stringify({
				error: error.message
			}))

		}

		response.end()
	})

	httpServer.listen(%s)`, port)

	file := os.TemporaryFile("server.js")
	err = os.Write(file, code, 0775)

	defer func() {
		if err != nil {
			os.Delete(file)
		}
	}()

	if err != nil {
		return err
	}

	// Run server
	cmd := exec.Command("node", file)
	err = cmd.Start()

	if err != nil {
		return err
	}

	go func() error {
		return cmd.Wait()
	}()

	address := "http://localhost:" + port

	// Wait service become online
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
		err := os.Delete(service.File)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute transpilation process
func (service *TranspilerService) Execute(config *TSConfig, file *compactor.File) error {

	data := struct {
		Config   *TSConfig `json:"config"`
		Source   string    `json:"source"`
		Filename string    `json:"filename"`
		Relative string    `json:"relative"`
	}{
		Config:   config,
		Filename: file.File,
		Relative: os.Relative(os.Dir(file.Destination), file.Path),
		Source:   file.Content,
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

	err = os.Write(file.Destination, result.Output, file.Permission)

	if err != nil {
		return err
	}

	if result.SourceMap != "" {
		err := os.Write(file.Destination+".map", result.SourceMap, file.Permission)

		if err != nil {
			return err
		}
	}

	return nil
}
