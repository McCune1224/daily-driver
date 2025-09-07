
.PHONY: run css templ build

run:
	air -c ./.air.toml

css:
	tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --watch

templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v


build:
	tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify
	templ generate
