FROM golang:1.24.4-alpine AS build

WORKDIR /build

COPY . .

RUN go mod download && go mod verify

WORKDIR /build/cmd/grpc

RUN go build -o post-service-grpc .

WORKDIR /build/cmd/bff

RUN go build -o post-service-bff .

FROM alpine

WORKDIR /app

COPY --from=build /build/.env.prod /app/.env

COPY --from=build /build/docs /app/docs
COPY --from=build /build/cmd/grpc/post-service-grpc /app/post-service-grpc
COPY --from=build /build/cmd/bff/post-service-bff /app/post-service-bff

RUN chmod +x /app/post-service-grpc /app/post-service-bff

EXPOSE 50051
EXPOSE 5001

CMD ["/bin/sh", "-c", "/app/post-service-grpc & /app/post-service-bff"]