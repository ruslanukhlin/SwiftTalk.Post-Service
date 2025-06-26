package main

// @title SwiftTalk Post Service API
// @version 1.0
// @description API сервиса постов для платформы SwiftTalk
// @host localhost:5001
// @BasePath /post-service/
import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/ruslanukhlin/SwiftTalk.post-service/docs"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/bff"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/post"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	server := fiber.New()

	server.Get("/docs/docs.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger-v3.json")
	})

	// OpenAPI 3.0 UI
	server.Get("/docs/*", swagger.New(swagger.Config{
		URL:          "docs.json",
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "SwiftTalk Post Service API (OpenAPI 3.0)",
	}))

	conn, err := grpc.NewClient(":"+cfg.PortGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}

	postClient := pb.NewPostServiceClient(conn)
	postService := bff.NewPostService(postClient)
	handler := bff.NewHandler(postService)

	bff.RegisterRoutes(server, handler)

	go func() {
		if err := server.Listen(":" + cfg.PortHttp); err != nil {
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

	if err := server.Shutdown(); err != nil {
		log.Printf("Ошибка graceful shutdown: %v", err)
	}
}
