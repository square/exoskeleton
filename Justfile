test:
    go test ./...

fmt:
    gofmt -w .
    goimports -w .

lint:
    gofmt -l -d .
    goimports -l -d .
    go vet ./...
