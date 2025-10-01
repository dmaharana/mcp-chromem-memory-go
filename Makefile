BINARY_NAME=memory-server

# Default build for the current OS
all: build

# Build for the current OS
build:
	go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Build for Linux
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BINARY_NAME)-linux ./cmd/$(BINARY_NAME)

# Build for Windows
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(BINARY_NAME)-windows.exe ./cmd/$(BINARY_NAME)

# clean up binaries
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux $(BINARY_NAME)-windows.exe

.PHONY: all build linux windows clean