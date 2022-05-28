package typescript

import (
	"encoding/json"
	"strings"

	"github.com/mateussouzaweb/compactor/os"
)

// Transpiler struct
type Transpiler struct {
	Config      *TSConfig
	File        string
	Content     string
	Destination string
}

// RunTranspiler will
func RunTranspiler(data Transpiler) error {

	runFile := os.TemporaryFile("transpile.js")
	sourceFile := os.TemporaryFile("source.ts")

	defer os.Delete(runFile)
	defer os.Delete(sourceFile)

	err := os.Write(sourceFile, data.Content, 0775)

	if err != nil {
		return err
	}

	config, err := json.Marshal(data.Config)

	if err != nil {
		return err
	}

	runCode := `
	const { execSync } = require("child_process")
	const root = execSync("npm root -g").toString().trim()
	const ts = require(root + "/typescript")
	const fs = require("fs")

	const config = JSON.parse('{CONFIG}')
		  config.fileName = '{RELATIVE}'

	const source = fs.readFileSync('{SOURCE}', {
		encoding: 'utf-8'
	})

	const result = ts.transpileModule(source, config)
	const output = result.outputText
	const sourceMap = result.sourceMapText.replace(
		'"sources":["{FILENAME}"]',
		'"sources":["{RELATIVE}"]'
	)

	fs.writeFileSync('{DESTINATION}', output)

	if (sourceMap) {
		fs.writeFileSync('{DESTINATION}.map', sourceMap)
	}
	`

	fileName := os.File(data.File)
	relative := os.Relative(os.Dir(data.Destination), data.File)

	runCode = strings.Replace(runCode, "{CONFIG}", string(config), -1)
	runCode = strings.Replace(runCode, "{SOURCE}", sourceFile, -1)
	runCode = strings.Replace(runCode, "{FILENAME}", fileName, -1)
	runCode = strings.Replace(runCode, "{RELATIVE}", relative, -1)
	runCode = strings.Replace(runCode, "{DESTINATION}", data.Destination, -1)

	err = os.Write(runFile, runCode, 0775)

	if err != nil {
		return err
	}

	// Run transpiler
	_, err = os.Exec(
		"node",
		runFile,
	)

	return err
}
