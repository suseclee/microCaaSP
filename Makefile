GO ?= GO111MODULE=on go

GOBINPATH     := $(shell $(GO) env GOPATH)/bin
#GOPATH="$HOME/go"

LN = ln
RM = rm


.PHONY: all
all: install

.PHONY: install
install:
	$(GO) install ./cmd/...

.PHONY: build
build:
	$(GO) build ./cmd/...

.PHONY: clean
clean:
	$(GO) clean -i ./...
	$(RM) -f ./microCaaSP
	$(RM) -f $(GOBINPATH)/microCaaSP

