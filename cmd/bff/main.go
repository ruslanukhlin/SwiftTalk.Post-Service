package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	_ "github.com/ruslanukhlin/SwiftTalk.post-service/docs"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/bff"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/post"
)

// @title SwiftTalk Post Service API
// @version 1.0
// @description API сервиса постов для платформы SwiftTalk
// @host localhost:5001
// @BasePath /postService/
func main() {
	cfg := config.LoadConfigFromEnv()

	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.FiberWrapHandler())

	conn, err := grpc.NewClient(":"+cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}

	postClient := pb.NewPostServiceClient(conn)
	postService := bff.NewPostService(postClient)
	handler := bff.NewHandler(postService)

	bff.RegisterRoutes(app, handler)

	go func() {
		if err := app.Listen(":" + cfg.PortHttp); err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("Ошибка при закрытии gRPC соединения: %v", err)
			}
		}()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := app.Shutdown(); err != nil {
		log.Printf("Ошибка graceful shutdown: %v", err)
	}
}
