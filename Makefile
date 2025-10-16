LOCAL_BIN:=$(CURDIR)/bin
LINTVER=v1.64.8

.PHONY:lint
lint:
	$(LOCAL_BIN)/golangci-lint run

.PHONY:install-deps
install-deps:
	@echo "Installing dependencies..."
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTVER)