.PHONY: build run get-dependencies

build: get-dependencies
	$(info Building Linux, Mac and Windows binaries)
	mkdir -p ./out/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/mvazquezc/termin8/pkg/version.gitCommit=$(shell git rev-parse HEAD)' -X 'github.com/mvazquezc/termin8/pkg/version.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/mvazquezc/termin8/pkg/version.version=$(shell git branch --show-current)'" -o ./out/termin8-linux-amd64 cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X 'github.com/mvazquezc/termin8/pkg/version.gitCommit=$(shell git rev-parse HEAD)' -X 'github.com/mvazquezc/termin8/pkg/version.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/mvazquezc/termin8/pkg/version.version=$(shell git branch --show-current)'" -o ./out/termin8-linux-arm64 cmd/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/mvazquezc/termin8/pkg/version.gitCommit=$(shell git rev-parse HEAD)' -X 'github.com/mvazquezc/termin8/pkg/version.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/mvazquezc/termin8/pkg/version.version=$(shell git branch --show-current)'" -o ./out/termin8-darwin-amd64 cmd/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'github.com/mvazquezc/termin8/pkg/version.gitCommit=$(shell git rev-parse HEAD)' -X 'github.com/mvazquezc/termin8/pkg/version.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/mvazquezc/termin8/pkg/version.version=$(shell git branch --show-current)'" -o ./out/termin8-darwin-arm64 cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/mvazquezc/termin8/pkg/version.gitCommit=$(shell git rev-parse HEAD)' -X 'github.com/mvazquezc/termin8/pkg/version.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/mvazquezc/termin8/pkg/version.version=$(shell git branch --show-current)'" -o ./out/termin8-windows-amd64.exe cmd/main.go
run: get-dependencies
	go run cmd/main.go
get-dependencies:
	$(info Downloading dependencies)
	go mod download