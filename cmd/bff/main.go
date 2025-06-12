package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/bff"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	app := fiber.New()

	conn, err := grpc.NewClient(":" + cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	postClient := pb.NewPostServiceClient(conn)
	postService := bff.NewPostService(postClient)
	handler := bff.NewHandler(postService)

	bff.RegisterRoutes(app, handler)

	if err := app.Listen(":" + cfg.PortHttp); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}