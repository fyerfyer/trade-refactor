package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
	"github.com/fyerfyer/trade-refactor/order/internal/application/service"
)

type Adapter struct {
	service *service.OrderService
	port    int
	pb.UnimplementedOrderServer
}

func NewAdapter(service *service.OrderService, port int) *Adapter {
	return &Adapter{
		service: service,
		port:    port,
	}
}

func (a *Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d:%v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServer(grpcServer, a)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d:%v", a.port, err)
	}
}
