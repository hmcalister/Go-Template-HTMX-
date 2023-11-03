# Standalone Tailwind CLI

First, [download the most recent Tailwind CLI version](https://github.com/tailwindlabs/tailwindcss/releases/latest). Then, [following the official docs](https://tailwindcss.com/blog/standalone-cli), move the tailwind file and make it executable. Initialize the `tailwind.config.js` file:

```bash
./tailwindcss init
```

Add the path to the template files to the `content` array in `tailwind.config.js`.

Create an `input.css` file and fill it with whatever tailwind classes required, e.g.:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Finally, start a watcher:

```bash
./tailwindcss -i input.css -o output.css --watch
```

To create a minified css file, use:

```bash
./tailwindcss -i input.css -o output.css --minify
```

Include in your templates the `output.css` file:

```html
<link rel='stylesheet' type='text/css' media='screen' href='css/output.css'>
```