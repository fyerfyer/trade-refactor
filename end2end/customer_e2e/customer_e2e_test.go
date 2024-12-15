package customere2e

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

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/customer"
)

type CustomerTestSuite struct {
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

func (c *CustomerTestSuite) SetupSuite() {
	if err := runDockerCompose("up", "--build", "-d"); err != nil {
		log.Fatalf("failed to run docker-compose up: %v", err)
	}
}

func (c *CustomerTestSuite) TearDownSuite() {
	if err := runDockerCompose("down"); err != nil {
		log.Printf("failed to stop docker-compose: %v", err)
	}
}

func (o *CustomerTestSuite) TestSubmitOrderService() {
	containerIP, err := getContainerIP("test_customer")
	o.Require().NoError(err, "failed to get container IP")

	connStr := fmt.Sprintf("%s:8083", containerIP)
	fmt.Printf("Connecting to: %s\n", connStr)

	// block the connect until success
	conn, err := grpc.Dial(connStr, grpc.WithInsecure(), grpc.WithBlock())
	o.Require().NoError(err, "did not connect")
	defer conn.Close()

	client := pb.NewCustomerClient(conn)

	// create the customer first
	createReq, createErr := client.CreateCustomer(context.Background(),
		&pb.CreateCustomerRequest{
			CustomerName: "Joe",
		})

	o.Nil(createErr)
	o.Equal(createReq.GetCustomerId(), uint64(1))
	o.Equal(createReq.GetSuccess(), true)

	// set the balance
	_, setErr := client.StoreBalance(context.Background(),
		&pb.StoreBalanceRequest{
			CustomerName: "Joe",
			Balance:      200,
		})

	o.Nil(setErr)

	reqOrder := &pb.SubmitOrderRequest{
		CustomerName: "Joe",
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

	resSubmit, errSubmit := client.SubmitOrder(context.Background(), reqOrder)
	o.Require().NoError(errSubmit)
	o.Require().NotNil(resSubmit)
	o.Equal(true, resSubmit.Success)
	o.Equal("successfully process items", resSubmit.Message)
}

func TestMain(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}
