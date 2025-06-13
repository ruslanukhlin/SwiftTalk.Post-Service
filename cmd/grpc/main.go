package main

import (
	"context"
	"log"
	"net"

	s3 "github.com/ruslanukhlin/SwiftTalk.common/core/s3"
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
	s3Client, err := s3.NewS3(context.Background(), cfg.S3.Bucket)
	if err != nil {
		log.Fatalf("Ошибка инициализации s3: %v", err)
	}
	postApp := application.NewPostApp(postDb, s3Client, cfg)

	runGRPCServer(postApp, cfg.PortGrpc)
}

func runGRPCServer(postApp *application.PostApp, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
	defer lis.Close()

	postGRPCHandler := postGRPC.NewPostGRPCHandler(postApp)
	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, postGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка grpc сервера: %v", err)
	}
}