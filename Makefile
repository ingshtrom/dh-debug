@PHONY: run
run: 
	@go generate ./... -n -v -x

	@export GOOS=linux
	@export CGO_ENABLED=0

	@go build -o main main.go
