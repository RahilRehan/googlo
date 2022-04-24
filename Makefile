CMD=up -d
all:
	docker-compose -f tools/docker-compose.yml $(CMD)

test:
	go clean -testcache
	go test ./...

.PHONY: all