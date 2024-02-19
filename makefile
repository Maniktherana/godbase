rs:
	go run ./cmd/redis-server/main.go

build:
	go build -o ./bin/redis/redis-server ./cmd/redis-server/main.go

PHONY: rs build



