package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
)

type OrderTestSuite struct {
	suite.Suite
}

func runDockerCompose(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = "../resource"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getContainerIP(containerName string) (string, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", containerName)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	ip := strings.TrimSpace(out.String())
	if ip == "" {
		return "", fmt.Errorf("could not find IP address for container %s", containerName)
	}
	return ip, nil
}

func (o *OrderTestSuite) SetupSuite() {
	if err := runDockerCompose("up", "--build", "-d"); err != nil {
		log.Fatalf("failed to run docker-compose up: %v", err)
	}
}

func (o *OrderTestSuite) TearDownSuite() {
	if err := runDockerCompose("down"); err != nil {
		log.Printf("failed to stop docker-compose: %v", err)
	}
}

func (o *OrderTestSuite) TestOrderService() {
	containerIP, err := getContainerIP("test_order")
	o.Require().NoError(err, "failed to get container IP")

	connStr := fmt.Sprintf("%s:8082", containerIP)
	fmt.Printf("Connecting to: %s\n", connStr)

	// block the connect until success
	conn, err := grpc.Dial(connStr, grpc.WithInsecure(), grpc.WithBlock())
	o.Require().NoError(err, "did not connect")
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
	o.Require().NoError(errItem)
	o.Require().NotNil(resItem)
	o.Equal("successfully process items", resItem.Message)

	resOrder, errOrder := client.GetOrder(context.Background(), &pb.GetOrderRequest{
		CustomerId: 1,
		Status:     "success",
	})
	o.Require().NoError(errOrder)
	o.Require().NotNil(resOrder)
	o.Equal(uint64(1), resOrder.GetOrder().OrderId)
	log.Printf("res: %+v\n", resOrder)
	log.Printf("items: %v", resOrder.Order.GetOrderItems())
}

func TestMain(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
