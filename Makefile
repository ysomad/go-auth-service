.DEFAULT_GOAL := build

# Go

build:
	go build -v ./cmd/server

test:
	go test -v -race -timeout 30s ./...


# Docker

dockerBuild:
	docker-compose -f ./build/docker-compose.yml build

dockerUp:
	docker-compose -f ./build/docker-compose.yml up -d

dockerDown:
	docker-compose -f ./build/docker-compose.yml down

dockerLogs:
	docker-compose -f ./build/docker-compose.yml logs

dockerBuildUp: dockerDown dockerBuild dockerUp


.PHONY: build test