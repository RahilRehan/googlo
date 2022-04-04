CMD=up -d
all:
	docker-compose -f tools/docker-compose.yml $(CMD)
	
.PHONY: all