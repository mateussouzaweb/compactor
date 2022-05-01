package typescript

import (
	"encoding/json"
	"strings"

	"github.com/mateussouzaweb/compactor/os"
)

// Transpiler struct
type Transpiler struct {
	Content     string
	Options     *TSConfig
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

	theOptions, err := json.Marshal(data.Options)

	if err != nil {
		return err
	}

	runCode := `
	const { execSync } = require("child_process")
	const root = execSync("npm root -g").toString().trim()
	const ts = require(root + "/typescript")
	const fs = require("fs")

	const options = JSON.parse('{OPTIONS}')
		  options.fileName = '{DESTINATION}'

	const source = fs.readFileSync('{SOURCE}', {
		encoding: 'utf-8'
	})

	const result = ts.transpileModule(source, options)
	const output = result.outputText
	const sourceMap = result.sourceMapText

	fs.writeFileSync('{OUTPUT}', output)

	if (sourceMap) {
		fs.writeFileSync('{OUTPUT}.map', sourceMap)
	}
	`

	runCode = strings.Replace(runCode, "{OPTIONS}", string(theOptions), 1)
	runCode = strings.Replace(runCode, "{DESTINATION}", os.File(data.Destination), 1)
	runCode = strings.Replace(runCode, "{SOURCE}", sourceFile, 1)
	runCode = strings.Replace(runCode, "{OUTPUT}", data.Destination, -1)

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
