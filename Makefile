VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
BINARY := ope

.PHONY: build build-macos build-windows build-linux icons test install clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

build-macos: build-macos-universal
	./scripts/make-app-bundle.sh dist/ope-darwin-universal $(VERSION) dist

build-macos-universal:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/ope-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ope-darwin-arm64 .
	lipo -create -output dist/ope-darwin-universal dist/ope-darwin-amd64 dist/ope-darwin-arm64

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/ope-windows-amd64.exe .

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ope-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ope-linux-arm64 .

icons:
	./scripts/make-icons.sh

test:
	go test ./... -v

install: build
	./$(BINARY) install

clean:
	rm -f $(BINARY)
	rm -rf dist/
