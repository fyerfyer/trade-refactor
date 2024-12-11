package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderRepoTestSuite struct {
	suite.Suite
	DBSourceURL string
}

func (o *OrderRepoTestSuite) SetupSuite() {
	ctx := context.Background()
	port := nat.Port("3306/tcp")

	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf(
			"root:110119abc@tcp(%s:%s)/order?charset=utf8&parseTime=True&loc=Local",
			host, port.Port(),
		)
	}

	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{string(port)},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "110119abc",
			"MYSQL_DATABASE":      "order",
		},
		WaitingFor: wait.ForSQL(port, "mysql", dbURL).
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(500 * time.Millisecond),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	o.Require().NoError(err, "failed to start MySQL container")

	host, err := mysqlContainer.Host(ctx)
	o.Require().NoError(err, "failed to get container host")

	mappedPort, err := mysqlContainer.MappedPort(ctx, port)
	o.Require().NoError(err, "failed to get container mapped port")

	o.DBSourceURL = dbURL(host, mappedPort)
}

func (o *OrderRepoTestSuite) Test_Save_Order() {
	gormRepo, err := NewGormRepository(o.DBSourceURL)
	o.Require().NoError(err, "failed to initialize gorm database")
	o.Require().NotNil(gormRepo, "repository is nil")

	order := &domain.Order{}
	saveErr := gormRepo.Save(context.Background(), order)
	o.NoError(saveErr, "failed to save order")
}

func (o *OrderRepoTestSuite) Test_Update_Order() {
	gormRepo, err := NewGormRepository(o.DBSourceURL)
	o.Require().NoError(err, "failed to initialize gorm database")
	o.Require().NotNil(gormRepo, "repository is nil")

	order := &domain.Order{
		ID:         1,
		CustomerID: 12345,
		Status:     "paid",
		Items: []domain.OrderItem{
			{ProductCode: "P001", UnitPrice: 10.0, Quantity: 2},
		},
		CreatedAt: time.Now(),
	}

	saveErr := gormRepo.Save(context.Background(), order)
	o.Require().NoError(saveErr, "failed to save order")

	order.Status = "unpaid"
	updateErr := gormRepo.Update(context.Background(), order)
	o.NoError(updateErr, "failed to update order")

	updatedOrder, getErr := gormRepo.GetUnpaidOrder(context.Background(), order.ID)
	o.NoError(getErr, "failed to get updated order")
	o.Equal("unpaid", updatedOrder.Status, "order status was not updated correctly")
}

func (o *OrderRepoTestSuite) Test_Delete_Order() {
	gormRepo, err := NewGormRepository(o.DBSourceURL)
	o.Require().NoError(err, "failed to initialize gorm database")
	o.Require().NotNil(gormRepo, "repository is nil")

	order := &domain.Order{
		CustomerID: 12345,
		Status:     "unpaid",
		Items: []domain.OrderItem{
			{ProductCode: "P001", UnitPrice: 10.0, Quantity: 2},
		},
	}

	saveErr := gormRepo.Save(context.Background(), order)
	o.Require().NoError(saveErr, "failed to save order")

	deleteErr := gormRepo.Delete(context.Background(), order.ID)
	o.NoError(deleteErr, "failed to delete order")

	_, getErr := gormRepo.GetUnpaidOrder(context.Background(), order.ID)
	o.Error(getErr, "order was not deleted correctly")
}

func (o *OrderRepoTestSuite) Test_Get_Unpaid_Orders() {
	gormRepo, err := NewGormRepository(o.DBSourceURL)
	o.Require().NoError(err, "failed to initialize gorm database")
	o.Require().NotNil(gormRepo, "repository is nil")

	orders := []domain.Order{
		{
			CustomerID: 12345,
			Status:     "unpaid",
			Items: []domain.OrderItem{
				{ProductCode: "P001", UnitPrice: 10.0, Quantity: 2},
			},
		},
		{
			CustomerID: 12345,
			Status:     "unpaid",
			Items: []domain.OrderItem{
				{ProductCode: "P002", UnitPrice: 20.0, Quantity: 1},
			},
		},
	}

	for _, order := range orders {
		saveErr := gormRepo.Save(context.Background(), &order)
		o.Require().NoError(saveErr, "failed to save order")
	}

	resultOrders, err := gormRepo.GetUnpaidOrders(context.Background(), 12345)
	o.NoError(err, "failed to retrieve unpaid orders")
	o.Equal(len(resultOrders), 2, "incorrect number of unpaid orders")
}

func (o *OrderRepoTestSuite) Test_Get_Unpaid_Order() {
	gormRepo, err := NewGormRepository(o.DBSourceURL)
	o.Require().NoError(err, "failed to initialize gorm database")
	o.Require().NotNil(gormRepo, "repository is nil")

	order := &domain.Order{
		CustomerID: 12345,
		Status:     "unpaid",
		Items: []domain.OrderItem{
			{ProductCode: "P001", UnitPrice: 10.0, Quantity: 2},
		},
	}

	saveErr := gormRepo.Save(context.Background(), order)
	o.Require().NoError(saveErr, "failed to save order")

	retrievedOrder, getErr := gormRepo.GetUnpaidOrder(context.Background(), order.ID)
	o.NoError(getErr, "failed to retrieve unpaid order")
	o.Equal(order.ID, retrievedOrder.ID, "incorrect order ID")
	o.Equal(order.Status, retrievedOrder.Status, "incorrect order status")
}

func TestOrderRepoSuite(t *testing.T) {
	suite.Run(t, new(OrderRepoTestSuite))
}
