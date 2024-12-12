package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/payment/config"
	"github.com/fyerfyer/trade-refactor/payment/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/payment/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/service"
)

func main() {
	dsn := config.GetDatabaseDSN()
	paymentRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init payment database: %v", err)
	}
	log.Println("successfully set up database connection")

	redisClient := redis.NewRedisClient(config.GetRedisAddr(), "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	paymentService := service.NewService(paymentRepo, redisClient)
	port := config.GetApplicationPort()
	grpcAdapter := grpc.NewAdapter(paymentService, port)
	log.Printf("payment grpc server is running on port %v ...", port)
	grpcAdapter.Run()
}
