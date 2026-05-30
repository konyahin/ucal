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

clean:
    rm -rf dist vscode-plugin/out
    rm -f ucal benchmark.txt coverage.out vscode-plugin/ucal-*.vsix

release version:
    rm -rf dist
    mkdir -p dist
    just _build linux   amd64 ucal-{{version}}-linux-amd64       ""
    just _build linux   arm64 ucal-{{version}}-linux-arm64       ""
    just _build darwin  amd64 ucal-{{version}}-darwin-amd64      ""
    just _build darwin  arm64 ucal-{{version}}-darwin-arm64      ""
    just _build windows amd64 ucal-{{version}}-windows-amd64     .exe
    just _package_vsix {{version}}
    cd dist && shasum -a 256 * > SHA256SUMS
    ls -lh dist

_build os arch name ext:
    GOOS={{os}} GOARCH={{arch}} go build -trimpath -ldflags="-s -w" -o dist/{{name}}{{ext}} .

_package_vsix version:
    cd vscode-plugin && npm ci
    cd vscode-plugin && npm run compile
    cd vscode-plugin && vsce package --no-update-package-json --out ../dist/ucal-{{version}}.vsix $(echo "{{version}}" | sed 's/^v//')
