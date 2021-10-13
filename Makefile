.PHONY: local-test
local-test:
	go run main.go

.PHONY: build
build: test
	CGO_ENABLED=0 go build -o ./bin/dh-debug main.go && \
	chmod +x ./bin/dh-debug

.PHONY: test
test: tidy
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

