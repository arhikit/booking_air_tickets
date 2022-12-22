RELATIVE_PATH=./cmd/app/
BINARY_NAME=main.exe

all: build test
 
build:
	go build -o ${BINARY_NAME} ${RELATIVE_PATH}main.go

gen-api:
	go generate ./...

lint: 
	golangci-lint run -c golangci.yml
 
test:
	go test -v ${RELATIVE_PATH}main.go
 
run:
	go run ${RELATIVE_PATH}main.go -c ${RELATIVE_PATH}config.yaml
 
clean:
	go clean
	rm ${BINARY_NAME}