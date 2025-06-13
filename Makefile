all: run

run:
	go run cmd/grpc/main.go & go run cmd/bff/main.go

swagger:
	swag init -g cmd/bff/main.go -d . -o ./docs

docker-up:
	docker compose --env-file .env.local up -d