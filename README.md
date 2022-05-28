# Compactor - Frontend compression without pain

*Compactor* is an efficient alternative to generate fully compressed HTML projects, including all JS, CSS images and any other resource files. You can use it to develop static websites or with your server side rendered website.

It was primarily designed to work with static websites where you don't have, don't need or don't want to use NodeJS ecosystem including NPM or Yarn - freedom is the key here.

The usage is very simple: you give the root folder of the project, and *compactor* builds the compressed version of the project for you with all possible optimizations. Start with no configuration needed.

----

## Features

- Written in Go Language in a single binary.
- Provided as package module for Go.
- Supports ignore, include and exclude rules.
- Optimizes HTML, CSS, SCSS, SASS, JavaScript, TypeScript, JSON and XML, ...
- Compiles SCSS/SASS to CSS.
- Compiles TypeScript to Javascript.
- Generates source-map for Javascript and CSS files.
- Automatically add hash ID to avoid caching in JS and CSS files: ``file.js`` -> ``file.485.js``
- Compress images in GIF, JPG, JPEG, PNG and SVG format.
- Automatically create WEBP copy from JPG, JPEG and PNG as progressive enhancement.
- Add support to HTML imports, so you can split the code and the system will automatically merge it on compilation.
- Develop mode for automation with file watcher and web server for live development.
- CLI flags to fine tuning control.
- Just works!

----

## RoadMap (In Development)

- Single output and merge for JSON, XML e SVG.
- Avif copy format from others images formats.
- Less, Stylus and CoffeeScript compilers.
- Support for VueJS, React, Svelte, ...
- PostCSS compilation.
- Font compression.
- More!

----

## CLI - Installation and Usage

Just run the command below to install compactor and dependencies:

```bash
curl https://mateussouzaweb.github.io/compactor/install.sh | bash
```

Done! To check command flags use:

```bash
compactor --help
```

To compress a project, run:

```bash
compactor \
    --source /path/to/source/ \
    --destination /path/to/destination/
```

To watch changes and live compress the project that is being rendered by other service, run:

```bash
compactor \
    --watch \
    --source /path/to/source/ \
    --destination /path/to/destination/
```

To run a complete dev environment for static projects, use the ``--develop`` option:

```bash
compactor \
    --develop true \
    --source /path/to/source/ \
    --destination /path/to/destination/
```

You can also run compactor with other modes and options. Check the available options with the ``--help`` flag.

----

## Usage with TypeScript - Required Options

To use TypeScript compilation, you must provide the ``tsconfig.json`` file with at least the following options. Please make sure that the ``--source`` CLI option are the same as the ``baseUrl`` value inside the config file:

```json
{
  "compilerOptions": {
    "baseUrl": "./src/",
    "isolatedModules": true
  }
}
```

That is it! Enjoy!
