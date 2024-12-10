package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/payment"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/order/internal/application/service"
)

func main() {
	dsn := "root:110119abc@tcp(127.0.0.1:3306)/microservice?charset=utf8&parseTime=true"
	orderRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init order database: %v", err)
	}
	log.Println("successfully set up database connection")

	paymentAdapter, err := payment.NewPaymentAdapter("localhost:50051")
	if err != nil {
		log.Fatalf("failed to set up payment grpc client: %v", err)
	}
	log.Println("Successfully dial to payment grpc server...")

	redisClient := redis.NewRedisClient("127.0.0.1:6379", "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	orderService := service.NewService(
		orderRepo,
		redisClient,
		paymentAdapter)

	grpcAdapter := grpc.NewAdapter(orderService, 50052)
	log.Println("order grpc server is running on port 50052...")
	grpcAdapter.Run()
}
