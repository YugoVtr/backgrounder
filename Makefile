generate-mock:
	go generate -v ./...

build:
	mkdir -p bin
	go build  -o bin/backgrounder cmd/main.go