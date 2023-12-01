# The following Makefile creates an abstraction to easily run the most common development commands.
# To run, just execute the command: `make <target>`. Excluding <target> builds the project by default.
BIN := simple-wd
# -trimpath makes Go stacktraces only print relative paths to files instead of full paths. Makes the
# terminal output cleaner.
BUILD_FLAGS ?= -trimpath
GO ?= go

# Builds the project. Default target when no targets are specified
.PHONY: build
build:
	$(GO) build $(BUILD_FLAGS)

# Build the project, then run it by calling the executable
.PHONY: build-run
build-run: build
	sudo ./$(BIN)

# Cleans build files, such as the executable
.PHONY: clean
clean:
	$(GO) clean -i ./...

# Formats all files accordingy to Go's standards and tidy the go.mod file
.PHONY: fmt
fmt:
	$(GO) fmt ./...
	$(GO) mod tidy -v

# Lints all files using both Go's built-in tool and a golangci-lint, a third-party linter
.PHONY: lint
lint:
	$(GO) vet ./...
	golangci-lint run ./...

# Run the project without compiling it
.PHONY: run
run:
	$(GO) run $(BUILD_FLAGS) .

# Run all tests for the project with a coverage report
.PHONY: test
test:
	$(GO) test -cover ./...
