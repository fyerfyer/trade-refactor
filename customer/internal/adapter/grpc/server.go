package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/customer"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/service"
	"google.golang.org/grpc"
)

type Adapter struct {
	service *service.CustomerService
	port    int
	pb.UnimplementedCustomerServer
}

func NewAdapter(service *service.CustomerService, port int) *Adapter {
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
	pb.RegisterCustomerServer(grpcServer, a)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d:%v", a.port, err)
	}
}

