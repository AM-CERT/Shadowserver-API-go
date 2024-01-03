.DEFAULT_GOAL := compile

COMMIT := $(shell git rev-parse --short HEAD)
PROJECT = Shadowserver-API-go
BIN_NAME = shadowserver-api-go
BUILD_DIR = ${PROJECT}-release-${COMMIT}
BUILD_FILE = "${BUILD_DIR}/${PROJECT}-release-${COMMIT}.zip"

compile:
	rm -rf $(PROJECT)-release-*
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BIN_NAME)-windows-amd64.exe github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/$(BIN_NAME)-windows-arm64.exe github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	GOOS=linux GOARCH=amd64   go build -o $(BUILD_DIR)/$(BIN_NAME)-linux-amd64       github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	GOOS=linux GOARCH=arm64   go build -o $(BUILD_DIR)/$(BIN_NAME)-linux-arm64       github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	GOOS=darwin GOARCH=amd64  go build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-amd64      github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	GOOS=darwin GOARCH=arm64  go build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-arm64      github.com/AM-CERT/$(PROJECT)/cmd/$(BIN_NAME)
	cp .env $(BUILD_DIR)/
	zip -r $(BUILD_FILE) $(BUILD_DIR)
