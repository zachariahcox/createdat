.PHONY: build clean test # build these rules even if files with the same name exist.
clean:
	rm -rf build

build:
	mkdir -p build
	go build -o build/$(BINARY_NAME) -v .

test:
	go test -v ./...
