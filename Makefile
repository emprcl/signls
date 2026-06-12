BIN                    := bin
BIN_NAME               := signls
GOLANG_BIN             := go
CGO_ENABLED            := 1
GOLANG_OS              := linux
GOLANG_ARCH            := amd64
GOLANG_BUILD_OPTS      += GOOS=$(GOLANG_OS)
GOLANG_BUILD_OPTS      += GOARCH=$(GOLANG_ARCH)
GOLANG_BUILD_OPTS      += CGO_ENABLED=$(CGO_ENABLED)
GOLANG_LINT            := $(BIN)/golangci-lint
GORELEASER             := github.com/goreleaser/goreleaser/v2@latest
ASEQDUMP_BIN           := aseqdump -p 14:0 | ts '[%H:%M:%.S]'
ASEQDUMP_NO_CLOCK_OPTS := | grep -v Clock

# GR_TARGET selects which goreleaser builds run (see .goreleaser.yaml). Default
# to the host so `make snapshot` builds a local binary.
HOST_OS                := $(shell go env GOOS)
HOST_ARCH              := $(shell go env GOARCH)
ifeq ($(HOST_OS),darwin)
GR_TARGET              ?= darwin
else
GR_TARGET              ?= $(HOST_OS)-$(HOST_ARCH)
endif

$(BIN):
	mkdir -p $(BIN)

$(GOLANG_LINT): $(BIN)
	GOBIN=$$(pwd)/$(BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: $(BIN)
	$(GOLANG_BUILD_OPTS) $(GOLANG_BIN) build -o $(BIN)/$(BIN_NAME)
	chmod +x $(BIN)/$(BIN_NAME)

# Build a local release snapshot (archives in dist/) for the host target.
snapshot:
	GR_TARGET=$(GR_TARGET) $(GOLANG_BIN) run $(GORELEASER) release --snapshot --clean

checks: $(GOLANG_LINT)
	$(GOLANG_LINT) run ./...

test:
	$(GOLANG_BIN) test ./...

bench:
	$(GOLANG_BIN) test -bench=. -benchmem ./...

monitor-midi:
	$(ASEQDUMP_BIN) $(ASEQDUMP_NO_CLOCK_OPTS)

monitor-midi-clock:
	$(ASEQDUMP_BIN)
