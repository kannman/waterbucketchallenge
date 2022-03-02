# keep binary deps installers here

SWAG_BIN=$(LOCAL_BIN)/swag
$(SWAG_BIN):
	$(info #Install swaggo $(SWAG_BIN))
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag


GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG=v1.44.2
GOLANGCI_PKG=github.com/golangci/golangci-lint/cmd/golangci-lint
$(GOLANGCI_BIN):
	$(info #Install golangci-lint $(GOLANGCI_TAG))
	cd /tmp && GOBIN=$(LOCAL_BIN) go get $(GOLANGCI_PKG)@$(GOLANGCI_TAG)
