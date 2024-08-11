.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate -watch -proxy="http://localhost:4000" -open-browser=false

.PHONY: air
air:
	air -c .air.toml

.PHONY: dev
dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./tmp ./cmd/main.go
	make air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Environment=production" -o ./bin ./cmd/main.go
.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test:
	  go test -v -timeout 30s ./...
