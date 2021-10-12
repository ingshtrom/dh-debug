.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./bin/dh-debug main.go && \
	chmod +x ./bin/dh-debug


.PHONY: test
test:
	go test ./...


