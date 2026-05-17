build:
    go build

run args="":
    go run . "{{ args }}"

test args="":
    go test ./... {{ args }}

benchmark:
    go test -run='^$' -bench=. -count=10 ./... > benchmark.txt
    cat benchmark.txt

benchstat: benchmark
    benchstat benchmark.txt

coverage:
    go test -cover ./...

coverage-html:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

check:
    go vet
    staticcheck
