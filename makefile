
BINARY = brc
# Default rule
all: build run clean

# Build
build:
    go build -o ${BINARY}

# Run
run:
    @echo "Usage: make run FILE=<file_path>"
    @if [ "${FILE}" = "" ]; then \
        echo "Error: FILE not specified. Use 'make run FILE=<file_path>' to run the project."; \
    else \
        ./${BINARY} ${FILE}; \
    fi

# Clean up
clean:
    if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

# Test
test:
    go test

# Help
help:
    @echo "make: build and run the project"
    @echo "make build: build the project"
    @echo "make run FILE=<file_path>: run the project with specified file"
    @echo "make clean: clean up the binary"
    @echo "make test: run tests"
    @echo "make help: show this message"

.PHONY: all build run clean test help
