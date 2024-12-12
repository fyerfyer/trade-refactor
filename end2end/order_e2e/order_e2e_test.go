package order_e2e

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
)

type OrderTestSuite struct {
	suite.Suite
	compose *tc.LocalDockerCompose
}

func (o *OrderTestSuite) SetupSuite() {
	log.Println("Starting SetupSuite...")
	composeFilePaths := []string{"../resource/docker-compose.yml"}
	log.Printf("Using docker-compose file paths: %v\n", composeFilePaths)
	identifier := strings.ToLower(uuid.New().String())
	log.Printf("Generated unique identifier: %s\n", identifier)
	compose := tc.NewLocalDockerCompose(composeFilePaths, identifier)
	o.compose = compose
	log.Println("Bringing up Docker compose stack...")
	err := compose.WithCommand([]string{"up", "-d"}).
		Invoke().
		Error
	if err != nil {
		log.Fatalf("Failed to run compose stack: %v\n", err)
	}
	log.Println("SetupSuite completed successfully.")
}

func (o *OrderTestSuite) Test_Process_Items() {
	log.Println("Starting Test_Process_Items...")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println("Dialing gRPC server at localhost:8082...")

	conn, err := grpc.Dial("localhost:8082", opts...)
	if err != nil {
		log.Fatalf("failed to connect payment service: %v", err)
	}
	defer func() {
		log.Println("Closing gRPC connection...")
		conn.Close()
	}()

	orderClient := pb.NewOrderClient(conn)
	log.Println("Creating order request...")

	clientRes, err := orderClient.ProcessItems(context.Background(),
		&pb.ProcessItemsRequest{
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
		})

	log.Printf("ProcessItems Response: %+v, Error: %v\n", clientRes, err)
	o.Nil(err, "the items processing should succeed")
	o.Equal(clientRes.GetMessage(), "successfully process items")

	log.Println("Fetching created order...")

	res, err := orderClient.GetOrder(context.Background(),
		&pb.GetOrderRequest{CustomerId: 1, Status: "unpaid"})
	log.Printf("GetOrder Response: %+v, Error: %v\n", res, err)

	o.Nil(err, "should get the order")

	if res != nil && res.Order != nil && len(res.Order.OrderItems) > 0 {
		orderItem := res.Order.OrderItems[0]
		log.Printf("Order Item: %+v\n", orderItem)
		o.Equal(float32(50.0), orderItem.UnitPrice)
		o.Equal(int32(1), orderItem.Quantity)
		o.Equal("juice", orderItem.ProductCode)
	} else {
		log.Println("Order items list is empty or nil")
	}

	log.Println("Test_Process_Items completed.")
}

func (o *OrderTestSuite) TearDownSuite() {
	log.Println("Tearing down Docker compose stack...")
	err := o.compose.
		WithCommand([]string{"down"}).
		Invoke().
		Error

	if err != nil {
		log.Fatalf("failed to tear down docker compose: %v", err)
	}
	log.Println("Docker compose stack torn down successfully.")
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
