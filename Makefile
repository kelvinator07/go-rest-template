include internal/config/.env

.PHONY: serve tidy test mock mig-up mig-down seed test-all postgres postgres-stop createdb dropdb redis redis-stop db-setup

serve:
	go run cmd/api/main.go
	
tidy:
	go mod tidy && go mod vendor

test:
	go run cmd/test/main.go

mock:
	@echo "Generating mocks for interface $(interface) in directory $(dir)..."
	mockery --name=$(interface) --dir=$(dir) --output=./internal/mocks
	cd ./internal/mocks && \
	mv $(interface).go $(filename).go

mig-up:
	go run cmd/migration/main.go -up

mig-down:
	go run cmd/migration/main.go -down

coverage:
	go test -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

seed:
	go run cmd/seed/main.go

test-all:
	go test -v -cover -short ./...

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:14-alpine

postgres-stop:
	docker stop postgres14 && docker rm postgres14

createdb:
	docker exec -it postgres14 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DATABASE}

dropdb:
	docker exec -it postgres14 dropdb ${POSTGRES_DATABASE}

redis:
	docker run --name redis -p ${REDIS_PORT}:${REDIS_PORT} -d redis:7-alpine --requirepass ${REDIS_PASS}

redis-stop:
	docker stop redis && docker rm redis

db-setup: 
	sh -c "make postgres; sleep 3; make createdb; sleep 3; make mig-up; sleep 3; make seed; sleep 3; make redis"
