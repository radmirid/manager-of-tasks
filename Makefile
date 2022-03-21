build:
	docker-compose build manager-of-tasks

run:
	docker-compose up manager-of-tasks

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up