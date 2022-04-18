ifeq ($(strip $(VERSION_STRING)),)
VERSION_STRING := $(shell git rev-parse --short HEAD)
endif

BINDIR    := $(CURDIR)/bin
PLATFORMS := linux/amd64/rk-Linux-x86_64 darwin/amd64/rk-Darwin-x86_64 windows/amd64/rk.exe linux/arm64/rk-Linux-arm64 darwin/arm64/rk-Darwin-arm64
BUILDCOMMAND := go build
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
label = $(word 3, $(temp))

UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
SHACOMMAND := shasum -a 256
else
SHACOMMAND := sha256sum
endif

.DEFAULT_GOAL := build

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
	go install github.com/uw-labs/strongbox@latest

.PHONY: test
test:
	go test -p 4 -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: integration-actions
integration-actions:
	go test g -tags="common actions" --count=1 -v -timeout 15m

.PHONY: dep
dep:
	go mod tidy

.PHONY: latest
latest:
	echo ${VERSION_STRING} > bin/latest

.PHONY: lint
lint:
	golangci-lint run

.PHONY: precommit
precommit:
	pre-commit run --all-files

.PHONY: build
build:
	$(BUILDCOMMAND) -o ${BINDIR}/previewd

.PHONY: release
build-all: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 $(BUILDCOMMAND) -o "bin/$(label)"
	$(SHACOMMAND) "bin/$(label)" > "bin/$(label).sha256"
