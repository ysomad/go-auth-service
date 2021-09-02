include .env.example
export

compose-up:
	docker-compose -f build/docker-compose.yml up --build -d postgres && docker-compose -f build/docker-compose.yml logs -f
.PHONY: compose-up

compose-down:
	docker-compose -f build/docker-compose.yml down --remove-orphans
.PHONY: compose-down

swag-v1:
	swag init -g internal/delivery/http/v1/router.go
.PHONY: swag-v1

run: swag-v1
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

docker-rm-volume:
	docker volume rm go-service-template_pg-data
.PHONY: docker-rm-volume

test:
	go test -v -cover -race ./internal/...
.PHONY: test

mock:
	mockery --all -r --case snake
.PHONY: mock

migrate-create:
	migrate create -ext sql -dir migrations $(name)
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down
