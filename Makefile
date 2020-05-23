MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

build:
	go build -v -o goor

lint:
	docker run --rm -v $(MAKEFILE_DIR):/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run --config=.github/.golangci.yml