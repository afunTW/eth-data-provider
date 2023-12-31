BINARY_NAME=eth_data_provider

install_dev:
	@go install github.com/swaggo/swag/cmd/swag@latest

gen_api_doc: install_dev
	@$(shell go env GOPATH)/bin/swag init
	@go generate ./...

build:
	@go mod download
	@go build -o ${BINARY_NAME} .

run: build
	./${BINARY_NAME}
