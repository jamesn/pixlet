GIT_COMMIT = $(shell git rev-list -1 HEAD)

ifeq ($(OS),Windows_NT)
	BINARY = pixlet.exe
	LDFLAGS = -ldflags="-s -extldflags=-static -X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
	TAGS = -tags timetzdata
else
	BINARY = pixlet
	LDFLAGS = -ldflags="-X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
	TAGS =
endif

all: build

test:
	go test $(TAGS) -v -cover ./...

test-coverage:
	go test $(TAGS) -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | tail -1

test-coverage-check:
	@go test $(TAGS) -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | tail -1
	@echo "Full report: coverage.html"

clean:
	rm -f $(BINARY)
	rm -rf ./build
	rm -rf ./out

bench:
	go test -benchmem -benchtime=20s -bench BenchmarkRunAndRender tidbyt.dev/pixlet/encode

build:
	go build $(LDFLAGS) $(TAGS) -o $(BINARY) tidbyt.dev/pixlet

embedfonts:
	go run render/gen/embedfonts.go
	gofmt -s -w ./

widgets:
	 go run runtime/gen/main.go
	 gofmt -s -w ./

release-macos: clean
	./scripts/release-macos.sh

release-linux: clean
	./scripts/release-linux.sh

release-windows: clean
	./scripts/release-windows.sh

install-buildifier:
	go install github.com/bazelbuild/buildtools/buildifier@latest

lint:
	@ buildifier --version >/dev/null 2>&1 || $(MAKE) install-buildifier
	buildifier -r ./

format: lint