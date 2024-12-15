package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/customer"
)

func getContainerIP(containerName string) (string, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", containerName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(out.String())
	if ip == "" {
		return "", fmt.Errorf("could not find IP address for container %s", containerName)
	}
	return ip, nil
}

func main() {
	containerName := "test_order"
	containerIP, err := getContainerIP(containerName)
	if err != nil {
		log.Fatalf("failed to get container IP: %v", err)
	}

	connStr := fmt.Sprintf("%s:8083", containerIP)
	fmt.Printf("Connecting to: %s\n", connStr)

	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCustomerClient(conn)

	_, createErr := client.CreateCustomer(context.Background(),
		&pb.CreateCustomerRequest{
			CustomerName: "Joe",
		})

	if createErr != nil {
		log.Fatalf("failed to create customer: %v", createErr)
	}

	_, setErr := client.StoreBalance(context.Background(),
		&pb.StoreBalanceRequest{
			CustomerName: "Joe",
			Balance:      200,
		})

	if setErr != nil {
		log.Fatalf("failed to set balance: %v", setErr)
	}

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
	if errSubmit != nil {
		log.Fatalf("failed to submit order: %v", errSubmit)
	}

	log.Printf("res: %+v\n", resSubmit)
}
