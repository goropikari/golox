BIN := mylang

build:
	go build -o $(BIN) ./cmd/...
