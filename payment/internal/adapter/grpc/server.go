package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/payment"
	"github.com/fyerfyer/trade-refactor/payment/internal/port"
	"google.golang.org/grpc"
)

type Adapter struct {
	service port.PaymentPort
	port int
	pb.UnimplementedPaymentServer
}

func NewAdapter(service port.PaymentPort, port int) *Adapter {
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
	pb.RegisterPaymentServer(grpcServer, a)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d:%v", a.port, err)
	}
}
