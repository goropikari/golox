BIN := tlps

build:
	go build -o $(BIN) ./cmd/...

test:
	go test -v
	bash ./test/test.sh

ast:
	go run tools/generate_ast.go ./ && go fmt .
