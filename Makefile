lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./main.go

test:
	go test ./pkg/bundler/ -v