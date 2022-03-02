LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=waterbucket

RUN_ARGS:=
include tools.mk

.PHONY: generate
generate: $(SWAG_BIN)
	$(SWAG_BIN) init --ot go,json -o api -g cmd/main.go

.PHONY: build
build: generate
	$(GOENV) CGO_ENABLED=0 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: lint
lint: $(GOLANGCI_BIN)
	$(GOENV) $(GOLANGCI_BIN) run --fix ./...

.PHONY: test
test:
	$(GOENV) go test -race -v ./...

.PHONY: run-docker-compose
run-docker-compose:
	 docker compose up --build

