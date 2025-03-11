default: fmt lint install generate

BINARY_NAME=terraform-provider-hashicups
DESTINATION=./bin/$(BINARY_NAME)

LINKER_FLAGS=-s -w -X 'github.com/hashicorp/terraform-provider-hashicups/version.ProviderVersion=local'

build:
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
