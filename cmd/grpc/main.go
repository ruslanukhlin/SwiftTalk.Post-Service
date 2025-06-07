package main

import (
	"log"
	"net"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/db/postgres"
	postGRPC "github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/grpc"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/gorm"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	if err := gorm.InitDB(config.DNS(cfg.Postgres)); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	if err := gorm.Migrate(cfg); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	postDb := postgres.NewPostgresMemoryRepository(gorm.DB)
	postApp := application.NewPostApp(postDb)

	runGRPCServer(postApp, cfg.Port)
}

func runGRPCServer(postApp *application.PostApp, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	postGRPCHandler := postGRPC.NewPostGRPCHandler(postApp)
	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, postGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка grpc сервера: %v", err)
	}
}