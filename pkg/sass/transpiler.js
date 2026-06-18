const { execSync } = require("child_process")
const root = execSync("npm root -g").toString().trim()
const sass = require(root + "/sass-embedded")
const http = require("http")
const url = require("url")
const port = process.env.PORT || 3000

const httpServer = http.createServer(async (request, response) => {

    try {

        const buffers = []
        for await (const chunk of request) {
            buffers.push(chunk)
        }

        const data = Buffer.concat(buffers).toString()
        const body = JSON.parse(data)

        const config = body.config
        const source = body.source
        const result = sass.compile(source, config)

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

httpServer.listen(port)