CMD=up -d
all:
	docker-compose -f tools/docker-compose.yml $(CMD)

sqlc-generate:
	docker run -it -v $(PWD):/usr kjconroy/sqlc -f /usr/tools/sqlc.yml generate

.PHONY: all