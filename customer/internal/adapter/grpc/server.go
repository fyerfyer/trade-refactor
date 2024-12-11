package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/customer"
	"github.com/fyerfyer/trade-refactor/customer/internal/port"
	"google.golang.org/grpc"
)

type Adapter struct {
	service port.CustomerPort
	port    int
	pb.UnimplementedCustomerServer
}

func NewAdapter(service port.CustomerPort, port int) *Adapter {
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
