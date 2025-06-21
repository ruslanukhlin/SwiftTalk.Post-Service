package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	s3 "github.com/ruslanukhlin/SwiftTalk.Common/core/s3"
	pbAuth "github.com/ruslanukhlin/SwiftTalk.Common/gen/auth"
	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	clientGRPC "github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/auth/client"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/db/postgres"
	postGRPC "github.com/ruslanukhlin/SwiftTalk.post-service/internal/infrastructure/post/grpc"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	// Проверяем доступность PostgreSQL
	log.Printf("Checking PostgreSQL connectivity...")
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", cfg.Postgres.Host, cfg.Postgres.Port), 5*time.Second)
	if err != nil {
		log.Printf("Warning: Cannot connect to PostgreSQL: %v", err)
	} else {
		conn.Close()
		log.Printf("PostgreSQL is reachable!")
	}

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
	authClient := clientGRPC.NewClientGRPC(getAuthClient(cfg))
	postApp := application.NewPostApp(postDb, s3Client, authClient)

	runGRPCServer(postApp, cfg.PortGrpc)
}

func runGRPCServer(postApp *application.PostApp, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("Ошибка при закрытии listener: %v", err)
		}
	}()

	postGRPCHandler := postGRPC.NewPostGRPCHandler(postApp)
	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, postGRPCHandler)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка grpc сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()
}

func getAuthClient(cfg *config.Config) pbAuth.AuthServiceClient {
	conn, err := grpc.NewClient(cfg.Auth.Host+":"+cfg.Auth.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка соединения с сервисом auth: %v", err)
	}
	return pbAuth.NewAuthServiceClient(conn)
}
