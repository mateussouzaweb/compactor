const { execSync } = require("child_process")
const root = execSync("npm root -g").toString().trim()
const ts = require(root + "/typescript")
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
        if (!data) {
            throw new Error("Empty request body")
        }

        const body = JSON.parse(data)
        const config = body.config || {}
        config.fileName = body.relative

        const source = body.source || ""
        const result = ts.transpileModule(source, config)

        const output = result.outputText ? result.outputText : ""
        const sourceMap = result.sourceMapText ? result.sourceMapText.replace(
            `"sources":["${body.filename}"]`,
            `"sources":["${body.relative}"]`
        ) : ""

        response.writeHead(200, { "Content-Type": "application/json" })
        response.end(JSON.stringify({
            success: true,
            output: output,
            sourceMap: sourceMap
        }))

    } catch (error){
        
        response.writeHead(400, { "Content-Type": "application/json" })
        response.end(JSON.stringify({
            success: false,
            message: error.message
        }))

    }
})

httpServer.listen(port, () => {
    console.log(`Transpiler running on port ${port}`)
})