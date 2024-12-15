package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/customer/config"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/order"
	"github.com/fyerfyer/trade-refactor/customer/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/service"
)

func main() {
	dsn := config.GetDatabaseDSN()
	customerRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init order database: %v", err)
	}
	log.Println("successfully set up database connection")

	redisClient := redis.NewRedisClient(config.GetRedisAddr(), "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	orderAdapter, err := order.NewOrderAdapter(config.GetOrderServiceAddr())
	log.Println("successfully dial to order grpc client...")

	customerService := service.NewService(
		customerRepo,
		redisClient,
		orderAdapter)

	grpcAdapter := grpc.NewAdapter(customerService, config.GetApplicationPort())
	log.Println("customer grpc server is running on port 50053...")
	grpcAdapter.Run()
}
