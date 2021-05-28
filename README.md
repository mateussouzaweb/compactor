# Compactor - Frontend compression without pain

*Compactor* is a simple alternative to generate a fully compressed HTML projects, including all JS, CSS and image files in the project. It was designed to work with static websites where you don't have, don't need or don't want to use NodeJS and their crap system but you can also use it with your server side rendered website.

The ideia is very simple: you give the root folder of the project, and *compactor* builds the compressed version of the project for you with all possible optimizations. Start with no configuration needed.

---

## Features

- Written in Go Language in a single binary.
- Provided as package module.
- Optimizes HTML, CSS, SCSS, SASS, JavaScript, TypeScript, JSON and XML, ...
- Compiles SCSS/SASS to CSS.
- Compiles TypeScript to Javascript.
- Compress images in GIF, JPG, JPEG, PNG and SVG format.
- Automatically create WEBP copy from JPG, JPEG and PNG as progressive enhancement.
- Watch mode for automation and live development.
- CLI flags to fine tuning control.
- Just works!

---

## RoadMap (In Development)

- File mapping and bundler feature
- HTML include feature for code splitting
- A way of create SourceMaps for HTML (maybe simple comment)
- Avif copy format from others images formats
- Less, Stylus and CoffeeScript compilers
- PostCSS compiler
- Font compression
- Chunk id to avoid caching in JS and CSS: [id].js -> 485.js
- More!

---

## CLI - Installation and Usage

For now, you first need to install dependencies:

```bash
sudo su
apt install -y jpegoptim libjpeg-progs optipng gifsicle webp nodejs npm

npm install -g html-minifier
npm install -g sass
npm install -g typescript
npm install -g uglify-js
npm install -g svgo
```

Then download the most recent binary file and make it executable:

```bash
GIT_URL=https://raw.githubusercontent.com/mateussouzaweb/compactor
sudo wget $GIT_URL/master/bin/compactor -O /usr/local/bin/compactor
sudo chmod +x /usr/local/bin/compactor
```

Done! To check command flags use:

```bash
compactor --help
```

To compress a project source into a destination, run:

```bash
compactor \
    --source /path/to/source/ \
    --destination /path/to/destination/
```

To watch live changes while developing a project, add the watch flag:

```bash
compactor \
    --watch
    --source /path/to/source/ \
    --destination /path/to/destination/
```

Enjoy!
