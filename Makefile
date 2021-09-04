BIN := mylang

build:
	go build -o $(BIN) ./cmd/...

test:
	go test -v
