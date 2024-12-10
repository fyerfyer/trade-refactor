package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/order"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/service"
)

func main() {
	dsn := "root:110119abc@tcp(127.0.0.1:3306)/microservice?charset=utf8&parseTime=true"
	customerRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init order database: %v", err)
	}
	log.Println("successfully set up database connection")

	redisClient := redis.NewRedisClient("127.0.0.1:6379", "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	orderAdapter, err := order.NewOrderAdapter("localhost:50052")
	log.Println("successfully dial to order grpc client...")

	customerService := service.NewService(
		customerRepo,
		redisClient,
		orderAdapter)

	grpcAdapter := grpc.NewAdapter(customerService, 50053)
	log.Println("customer grpc server is running on port 50053...")
	grpcAdapter.Run()
}
