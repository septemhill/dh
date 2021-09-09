.PHONY: build clean run run-containers clean-containers patch_sqls

POSTGRES_PASSWORD=simple_password
DOCKER=docker

all: build

build:
	@go build

run: build run-containers patch_sqls
	@REDIS_ADDR="localhost:6379" \
	REDIS_PASSWORD="" \
	POSTGRES_CONNECTION_INFO="host=localhost user=postgres dbname=postgres password=$(POSTGRES_PASSWORD) port=5432 sslmode=disable" \
	SHORT_URL_PORT="80" \
	./dh

patch_sqls:
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -h localhost -U postgres -f sqls/0001_create_short_url_table.sql

run-containers:
	@$(DOCKER) compose up -d 
	@sleep 2

clean-containers:
	@$(DOCKER) compose stop
	@$(DOCKER) compose rm -f

clean: clean-containers
	@rm -rf dh
