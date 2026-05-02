build:
    go build

test args:
    go test ./... {{ args }}

coverage:
    go test -cover ./...

coverage-html:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

check:
    go vet
    staticcheck
