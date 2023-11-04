dev:
	air --build.cmd "go build -o tmp/main" --build.exclude_dir "static/css" --build.send_interrupt True
	