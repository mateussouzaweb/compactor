# Compactor - Frontend compression without pain

*Compactor* is an efficient alternative to generate fully compressed HTML projects, including all JS, CSS images and any other resource files. You can use it to develop static websites or with your server side rendered website.

It was primarily designed to work with static websites where you don't have, don't need or don't want to use NodeJS ecosystem including NPM or Yarn - freedom is the key here.

The usage is very simple: you give the root folder of the project, and *compactor* builds the compressed version of the project for you with all possible optimizations. Start with no configuration needed.

----

## Features

- Written in Go as a single binary.
- Provided as a Go module.
- Supports ignore, include and exclude rules.
- Optimizes HTML, CSS, SCSS, SASS, JavaScript, TypeScript, JSON and XML, ...
- Compiles SCSS/SASS to CSS.
- Compiles TypeScript to JavaScript.
- Generates source maps for JavaScript and CSS files.
- Automatically adds a hash ID to avoid caching in JS and CSS files: ``file.js`` -> ``file.485.js``
- Compresses images in GIF, JPG/JPEG, PNG and SVG formats.
- Automatically creates a WEBP copy from JPG/JPEG and PNG as a progressive enhancement.
- Adds support for HTML imports, so you can split the code and the system will automatically merge it on compilation.
- Develop mode for automation with file watcher and web server for live development.
- CLI flags for fine‑tuning control.
- Just works!

----

## Roadmap (In Development)

- Single output and merge for JSON, XML and SVG.
- Add AVIF output generation from other image formats.
- Less, Stylus and CoffeeScript compilers.
- Support for VueJS, React, Svelte, ...
- PostCSS compilation.
- Font compression.
- More!

----

## Docker - Usage

You can run *Compactor* directly from a Docker container image. Simply pull the published image from the registry:

```bash
docker pull ghcr.io/mateussouzaweb/compactor:latest
```

Now, to check command flags use:

```bash
docker run --rm ghcr.io/mateussouzaweb/compactor:latest --help
```

To run the image against a local source directory and write output to a destination directory:

```bash
docker run --rm \
  -v "$PWD/src:/src" \
  -v "$PWD/dist:/dist" \
  ghcr.io/mateussouzaweb/compactor:latest \
  --source /src --destination /dist
```

To run in watch mode (rendered by other service):

```bash
docker run --rm \
  -v "$PWD/src:/src" \
  -v "$PWD/dist:/dist" \
  ghcr.io/mateussouzaweb/compactor:latest \
  --source /src --destination /dist --watch
```

To run in development mode (watch + local server):

```bash
docker run --rm \
  -v "$PWD/src:/src" \
  -v "$PWD/dist:/dist" \
  -p 5000:5000 \
  ghcr.io/mateussouzaweb/compactor:latest \
  --source /src --destination /dist \
  --server :5000 --develop true
```

----

## CLI - Installation and Usage

Just run the command below to install compactor and dependencies:

```bash
curl https://mateussouzaweb.github.io/compactor/install.sh | bash -
```

To check command flags use:

```bash
compactor --help
```

To compress a project, run:

```bash
compactor \
  --source src/ \
  --destination dist/
```

To watch changes and live compress the project that is being rendered by other service, run:

```bash
compactor \
  --watch \
  --source src/ \
  --destination dist/
```

To run a complete dev environment for static projects, use the ``--develop`` option:

```bash
compactor \
  --develop true \
  --server :5000 \
  --source src/ \
  --destination dist/
```

You can also run compactor with other modes and options. Check the available options with the ``--help`` flag.

----

## Usage with TypeScript - Required Options

To use TypeScript compilation, you must provide the ``tsconfig.json`` file with at least the following options. Please make sure the `--source` CLI option matches the `baseUrl` value inside the config file:

```json
{
  "compilerOptions": {
    "baseUrl": "./src/",
    "isolatedModules": true
  }
}
```

That is it! Enjoy!