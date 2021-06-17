lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

test:
	go test ./pkg/bundler/ -v