include .env.example
export

compose-up:
	docker-compose -f ./build/docker-compose.yml up --build -d postgres mongodb redis && docker-compose -f ./build/docker-compose.yml logs -f
.PHONY: compose-up

compose-down:
	docker-compose -f ./build/docker-compose.yml down --remove-orphans
.PHONY: compose-down

run:
	go mod tidy && go mod download && \
	GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

migrate-create:
	migrate create -ext sql -dir migrations $(name)
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down
