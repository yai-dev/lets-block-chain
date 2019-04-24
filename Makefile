GOCMD=go
GOTEST=${GOCMD} test
SRC=$(shell find . -type f -name '*_test.go')
TARGET=tool

.PHONY: lint test build

all: .PHONY

build: $(TARGET)

tool:
	$(GOCMD) build -o lbc cmd/tool.go

test:
	@mkdir ./db
	@$(GOTEST) -v -cover ./... 2>&1
	@rm -rf ./db

clean:
	@rm -f ./lbc
	@rm -rf ./db
