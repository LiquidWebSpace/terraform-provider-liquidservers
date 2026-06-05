BINARY=terraform-provider-liquidservers

.PHONY: build test tidy install-local

build:
	go build -o bin/$(BINARY) .

test:
	go test ./...

tidy:
	go mod tidy

install-local: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/liquidservers/liquidservers/0.1.0/linux_amd64
	cp bin/$(BINARY) ~/.terraform.d/plugins/registry.terraform.io/liquidservers/liquidservers/0.1.0/linux_amd64/$(BINARY)

