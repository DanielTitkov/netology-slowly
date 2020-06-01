ENTRY_PATH = ./cmd
TEST_PATH = ./test

.PHONY: run
run: tidy
	go run $(ENTRY_PATH)/main.go

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test: tidy
	go test $(TEST_PATH) -v

.PHONY: lint
lint: tidy
	@golangci-lint run