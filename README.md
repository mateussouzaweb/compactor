# Compactor - Frontend compression without pain

*Compactor* is a efficient alternative to generate fully compressed HTML projects, including all JS, CSS images and ony other resource files.

It was primary designed to work with static websites where you don't have, don't need or don't want to use NodeJS and their crap system but you can also use it with your server side rendered website.

The ideia is very simple: you give the root folder of the project, and *compactor* builds the compressed version of the project for you with all possible optimizations. Start with no configuration needed.

---

## Features

- Written in Go Language in a single binary.
- Provided as package module for Go.
- Support files mapping and bundler.
- Support ignore, include and exclude rules.
- Optimizes HTML, CSS, SCSS, SASS, JavaScript, TypeScript, JSON and XML, ...
- Compiles SCSS/SASS to CSS.
- Compiles TypeScript to Javascript.
- Generates source-map for Javascript and CSS files.
- Automatically add hash id to avoid caching in JS and CSS files: ``file.js`` -> ``file.485.js``
- Compress images in GIF, JPG, JPEG, PNG and SVG format.
- Automatically create WEBP copy from JPG, JPEG and PNG as progressive enhancement.
- Watch mode for automation and live development.
- CLI flags to fine tuning control.
- Just works!

---

## RoadMap (In Development)

- HTML include feature for code splitting.
- A way of create source-map for HTML (maybe simple comment).
- Single output from multiple files for SASS and TypeScript (these languages does not include it out of the box, you have to use a 'index' file).
- Single output and merge for JSON, XML e SVG.
- Avif copy format from others images formats.
- Less, Stylus and CoffeeScript compilers.
- Support for VueJS, React, Svelte, ...
- PostCSS compilation.
- Font compression.
- More!

---

## CLI - Installation and Usage

For now, you first need to install dependencies:

```bash
sudo apt install -y nodejs npm libjpeg-progs
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

---

## File Mappings and Bundles

If you are using custom file mappings or package bundle features, please make sure to require the correct file or path in our HTML files and related code that require such files. You should always use the final file path.

For example, if you have 2 scripts (``lib.js`` and ``events.js``) that are merged and placed in ``scripts.js`` by ``compactor``, set the reference of the file like bellow:

```html
<script src="scripts.js"></script>
```

The same applies to CSS, JS files or any other file that use others files as reference:

```css
/* CSS */
.bg {
    background-image: url('final-name.png')
}
```

```js
// JS
required('bundled.js')
```

That is it! Enjoy!
