.ONESHELL:

BINARY_NAME=fs_producer

# golang commands
fmt:
	go fmt ./...

lint: 
	golint ./...

test:
	go test -v -cover ./...

build:
	go build -o bin/${BINARY_NAME} ./cmd

buildexe:
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}.exe ./

run: build 

	./bin/${BINARY_NAME}

.PHONY: build_docker push_docker fmt lint test
