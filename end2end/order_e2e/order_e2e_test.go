package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
)

type OrderTestSuite struct {
	suite.Suite
	orderHost string
	orderPort string
}

func runDockerCompose(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = "../resource"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *OrderTestSuite) SetupSuite() {
	if err := runDockerCompose("up", "--build", "-d"); err != nil {
		log.Fatalf("failed to run docker-compose up: %v", err)
	}

	containerRequest := testcontainers.ContainerRequest{
		Name: "test_order",
	}

	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
			Reuse:            true,
		},
	)
	if err != nil {
		log.Fatalf("failed to connect to existing container: %v", err)
	}

	host, err := container.Host(context.Background())
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(context.Background(), "8082")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	o.orderHost = host
	o.orderPort = port.Port()
}

func (o *OrderTestSuite) TearDownSuite() {
	if err := runDockerCompose("down"); err != nil {
		log.Printf("failed to stop docker-compose: %v", err)
	}
}

func (o *OrderTestSuite) TestOrderService() {
	connStr := fmt.Sprintf("%s:%s", o.orderHost, o.orderPort)
	fmt.Printf("Connecting to: %s\n", connStr)

	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderClient(conn)

	reqOrder := &pb.ProcessItemsRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: 1,
			Balance:    200.0,
		},
		OrderItems: []*pb.OrderItem{
			{
				ProductCode: "juice",
				UnitPrice:   50.0,
				Quantity:    1,
			},
			{
				ProductCode: "bread",
				UnitPrice:   50.0,
				Quantity:    2,
			},
		},
	}

	resItem, errItem := client.ProcessItems(context.Background(), reqOrder)
	o.NoError(errItem)
	o.NotNil(resItem)
	o.Equal(resItem.Message, "successfully process items")

	resOrder, errOrder := client.GetOrder(context.Background(), &pb.GetOrderRequest{
		CustomerId: 1,
		Status:     "success",
	})

	o.NoError(errOrder)
	o.NotNil(resOrder)
	o.Equal(resOrder.GetOrder().OrderId, uint64(1))
	log.Printf("res: %+v\n", resOrder)
	log.Printf("items: %v", resOrder.Order.GetOrderItems())
}

func TestMain(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
