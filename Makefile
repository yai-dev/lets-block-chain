GOCMD=go
GOTEST=${GOCMD} test
SRC=$(shell find . -type f -name '*.go')

.PHONY: lint test

all: .PHONY

test:
	@$(GOTEST) -v -cover ./... 2>&1