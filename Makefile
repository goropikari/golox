BIN := tlps

build:
	go build -o $(BIN) ./cmd/...

test:
	go test -v

ast:
	go run tools/generate_ast.go ./ && go fmt .
