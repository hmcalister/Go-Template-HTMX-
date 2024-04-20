# GoHTMXTemplate

A quick investigation into using Go Templates with [HTMX](https://htmx.org/). This project also serves as a good template for developing applications with an HTML frontend using Go Templates and HTMX. We also include [TailwindCSS](https://tailwindcss.com/docs/installation) and [DaisyUI](https://daisyui.com) for easy styling. All files are embedded into the compiled binary in this example, but this can be changed if required.

## Setting Up This Project

### Using `gonew`

[GoNew](https://go.dev/blog/gonew) is a tool to create a new go project from a template repository. To use it with this repo, simply run:

```bash
gonew github.com/hmcalister/GoHTMXTemplate YOUR_PROJECT_NAME
cd YOUR_PROJECT_NAME
make init
```

And start developing!

### Manually

First, git clone this project and remove the `.git` directory:

```bash
mkdir NewProject; cd NewProject
git clone https://github.com/hmcalister/GoHTMXTemplate .
rm -rf .git
make init
```

(Make sure you rename the Go module in `go.mod` to whatever you need!)

Use `go run .` to run the example, which will serve the index webpage at `http://localhost:8080/`. Navigate to this address with your web-browser to see HTMX in action.

When ready, add your own application code to the `api` package, creating methods that will be called by HTTP requests to the server. It is easiest to have a series of structs that hold your application state and can be updated with simple function calls. Add any required api handlers to `main.go`. Take note of the example handlers for how HTML templates can be executed with application state to return HTML code, as HTMX requires.

Note that `/css/output.css` and `/htmx/htmx.js` are embedded and served statically. If your frontend changes drastically (or you change the routes to the front end) ensure your server also serves these static files in the correct place.

## CSS - TailwindCSS and DaisyUI

Both TailwindCSS and DaisyUI are employed make styling easy. This will require Node to be installed, however if you are okay with leaving out DaisyUI (and any other plugins) then it is possible to run TailwindCSS using the [standalone CLI](https://tailwindcss.com/blog/standalone-cli). 

To set up the CSS, navigate to `static/css` and run `npm install -D tailwindcss daisyui@latest`. The already present `tailwind.config.js` file should handle the remaining setup, although if any problems arise it may be best to follow the setup for both TailwindCSS and DaisyUI individually yourself. 

After this, from the `static/css` directory, run `npx tailwindcss -i input.css -o output.css --watch` to continually watch the HTML files for any changes, and update the `output.css` file when this occurs. To build the CSS file in a minified format for deployment, run instead `npx tailwindcss -i input.css -o output.css --minify`.

Of course, you are free to replace all of the TailwindCSS + DaisyUI implementation with your own CSS if you want. Just note that currently the project is hardcoded to embed and server `static/css/output.css`.

## Makefile + Requirements

A Makefile is supplied for ease of development. Running `make dev` will run both the TailwindCSS watcher and a watcher for the Go server. The watcher uses [Cosmtrek's Air](https://github.com/cosmtrek/air) project, which is a requirement to use the `make dev` command. Once running, any changes to the Go files or HTML templates should result in a rebuild and rerun of the server, although you will need to reload the webpage in your web-browser.

`make build` will minify the CSS and build the project. This does not require Air, but does assume you are using TailwindCSS and Node.