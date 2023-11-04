dev:
	cd static/css; npx tailwindcss -i input.css -o output.css --watch &
	air --build.cmd "go build -o tmp/main" --build.exclude_dir "static/css" --build.send_interrupt True

build:
	cd static/css; npx tailwindcss -i input.css -o output.css --minify
	go build .
