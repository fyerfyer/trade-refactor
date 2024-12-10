package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/payment/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/payment/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/service"
)

func main() {
	dsn := "root:110119abc@tcp(127.0.0.1:3306)/microservice?charset=utf8&parseTime=true"
	paymentRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init payment database: %v", err)
	}
	log.Println("successfully set up database connection")

	redisClient := redis.NewRedisClient("127.0.0.1:6379", "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	paymentService := service.NewService(paymentRepo, redisClient)
	grpcAdapter := grpc.NewAdapter(paymentService, 50051)
	log.Println("payment grpc server is running on port 50051...")
	grpcAdapter.Run()
}
