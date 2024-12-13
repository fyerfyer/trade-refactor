package main

import (
	"log"
	"time"

	redis "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/order/config"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/grpc"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/payment"
	"github.com/fyerfyer/trade-refactor/order/internal/adapter/repo"
	"github.com/fyerfyer/trade-refactor/order/internal/application/service"
	google "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	dsn := config.GetDatabaseDSN()
	orderRepo, err := repo.NewGormRepository(dsn)
	if err != nil {
		log.Fatalf("failed to init order database: %v", err)
	}
	log.Println("successfully set up database connection")

	paymentAdapter, err := payment.NewPaymentAdapter("localhost:8081")
	if err != nil {
		log.Fatalf("failed to set up payment grpc client: %v", err)
	}
	log.Println("Successfully dial to payment grpc server...")

	redisClient := redis.NewRedisClient(config.GetRedisAddr(), "", 10, 10, 3*time.Minute)
	log.Println("successfully set up redis connection")

	// health check
	srv := google.NewServer()
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	orderService := service.NewService(
		orderRepo,
		redisClient,
		paymentAdapter)

	grpcAdapter := grpc.NewAdapter(orderService, config.GetApplicationPort())
	log.Println("order grpc server is running on port 50052...")
	grpcAdapter.Run()
}
