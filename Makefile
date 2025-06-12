all: run

run:
	go run cmd/grpc/main.go & go run cmd/bff/main.go

docker-up:
	docker compose --env-file .env.local up -d