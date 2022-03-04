
.PHONY: run.server
run.server:
	CONFIG_NAME=$(env) go run $(CURDIR)/main.go server


.PHONY: migrate.create
migrate.create:
	goose -dir $(CURDIR)/build/migrations create $(m) sql

.PHONY: migrate.local.status
migrate.local.status:
	goose -dir $(CURDIR)/build/migrations mysql "test:test@tcp(localhost:3306)/cro_local?parseTime=true" status


.PHONY: migrate.local.up
migrate.local.up:
	goose -dir $(CURDIR)/build/migrations mysql "test:test@tcp(localhost:3306)/cro_local?parseTime=true" up


.PHONY: migrate.local.down
migrate.local.down:
	goose -dir $(CURDIR)/build/migrations mysql "test:test@tcp(localhost:3306)/cro_local?parseTime=true" down

.PHONY: run.local.mysql
run.local.mysql:
	docker run -d -p 3306:3306 \
    	-e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=cro_local \
		-e MYSQL_USER=test -e MYSQL_PASSWORD=test \
    	--name db mysql/mysql-server:5.7

.PHONY: run.local.redis
redis.local.run:
	docker run --name redis -d -p 6379:6379 redis

run.test.mysql:
	docker run -d -p 13306:3306 \
    	-e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=cro_test \
		-e MYSQL_USER=test -e MYSQL_PASSWORD=test \
    	--name test-mysql mysql/mysql-server:5.7

.PHONY: build.image
build.image:
	docker build -t vicxu/cro-test -f $(CURDIR)/build/docker/server.dockerfile .
	docker push vicxu/cro-test

.PHONY: compose.up
compose.up:
	> $(CURDIR)/log/output.log
	docker-compose -f $(CURDIR)/deployment/local/docker-compose.yaml --env-file $(CURDIR)/deployment/local/.env up

.PHONY: init.swag
init.swag:
	swag init -g ./cmd/server.go 

