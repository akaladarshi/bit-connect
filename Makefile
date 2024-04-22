.PHONY: build install connect

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME)

install: go.sum
	@echo "installing bit-connect binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/bit-connect main.go

connect:
	@echo "Connecting to bitcoin node (0.0.0.0:18444)..."
	@go run main.go connect