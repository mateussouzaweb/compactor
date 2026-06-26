const { execSync } = require("child_process")
const root = execSync("npm root -g").toString().trim()
const sass = require(root + "/sass-embedded")
const http = require("http")
const path = require("path")
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
        const file = body.file || {}
        const content = file.content || ""
        const config = body.config || {}
        config.loadPaths = [path.dirname(file.path)]

        const result = await sass.compileStringAsync(content, config);
        const output = result.css ? result.css.toString() : ""
        const sourceMap = result.sourceMap
            ? JSON.stringify(result.sourceMap).replace(
                `"sources":["${file.file}"]`,
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