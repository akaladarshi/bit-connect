.PHONY: build install

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME)

install: go.sum
	@echo "installing bit-connect binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/bit-connect main.go