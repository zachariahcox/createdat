GIT_TAG?= $(shell git describe --always --tags)
BIN = gosec
FMT_CMD = $(gofmt -s -l -w $(find . -type f -name '*.go' -not -path './vendor/*') | tee /dev/stderr)
IMAGE_REPO = securego
DATE_FMT=+%Y-%m-%d
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
BUILDFLAGS := "-w -s -X 'main.Version=$(GIT_TAG)' -X 'main.GitTag=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)'"
BUILD_DIR := $(CURDIR)/build

.PHONY: build clean test # build these rules even if files with the same name exist.
clean:
	rm -rf $(BUILD_DIR)

build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -tags=netgo -ldflags=$(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

test:
	CGO_ENABLED=0 go test -v ./...
