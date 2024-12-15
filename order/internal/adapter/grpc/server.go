package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
	"github.com/fyerfyer/trade-refactor/order/internal/port"
)

type Adapter struct {
	service port.OrderPort
	port    int
	pb.UnimplementedOrderServer
}

func NewAdapter(service port.OrderPort, port int) *Adapter {
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
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d:%v", a.port, err)
	}
}
