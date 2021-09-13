BIN := golox

build:
	go build -o $(BIN) ./cmd/...

test: build
	go test -v
	bash ./test/test.sh

ast:
	go run tools/generate_ast.go ./ && go fmt .
