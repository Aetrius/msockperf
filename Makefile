# Makefile

# Define the name of the output binary
BINARY_NAME=msockperf-ubuntu

# Set the target operating system and architecture
GOOS=linux
GOARCH=amd64

# List all source files
SOURCES=main.go msockperf.go

# Build the binary
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) $(SOURCES)

# Clean the project
clean:
	rm -f $(BINARY_NAME)
