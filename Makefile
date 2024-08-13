BIN                    := bin
BIN_NAME               := cykl
GOLANG_BIN             := go
CGO_ENABLED            := 1
GOLANG_OS              := linux
GOLANG_ARCH            := amd64
GOLANG_BUILD_OPTS      += GOOS=$(GOLANG_OS)
GOLANG_BUILD_OPTS      += GOARCH=$(GOLANG_ARCH)
GOLANG_BUILD_OPTS      += CGO_ENABLED=$(CGO_ENABLED)
GOLANG_LINT            := $(BIN)/golangci-lint
ASEQDUMP_BIN           := aseqdump -p 14:0 | ts '[%H:%M:%.S]'
ASEQDUMP_NO_CLOCK_OPTS := | grep -v Clock

$(BIN):
	mkdir -p $(BIN)

$(GOLANG_LINT): $(BIN)
	GOBIN=$$(pwd)/$(BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: $(BIN)
	$(GOLANG_BUILD_OPTS) $(GOLANG_BIN) build -o $(BIN)/$(BIN_NAME)
	chmod +x $(BIN)/$(BIN_NAME)

checks: $(GOLANG_LINT)
	$(GOLANG_LINT) run ./...

bench:
	$(GOLANG_BIN) test -bench=. -benchmem ./...

monitor-midi:
	$(ASEQDUMP_BIN) $(ASEQDUMP_NO_CLOCK_OPTS)

monitor-midi-clock:
	$(ASEQDUMP_BIN)
