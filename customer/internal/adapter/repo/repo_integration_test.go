package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	DBSourceURL string
}

func (c *CustomerRepoTestSuite) SetupSuite() {
	ctx := context.Background()
	port := nat.Port("3306/tcp")

	// Define database URL function
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf(
			"root:110119abc@tcp(%s:%s)/customer?charset=utf8&parseTime=True&loc=Local",
			host, port.Port(),
		)
	}

	// Define container request
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{string(port)},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "110119abc",
			"MYSQL_DATABASE":      "customer",
		},
		WaitingFor: wait.ForSQL(port, "mysql", dbURL).
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(500 * time.Millisecond),
	}

	// Start the MySQL container
	mysqlContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	c.Require().NoError(err, "failed to start MySQL container")

	// Get host, port from the container
	host, err := mysqlContainer.Host(ctx)
	c.Require().NoError(err, "failed to get container host")

	mappedPort, err := mysqlContainer.MappedPort(ctx, port)
	c.Require().NoError(err, "failed to get container mapped port")

	c.DBSourceURL = dbURL(host, mappedPort)
}

func (c *CustomerRepoTestSuite) TestSaveCustomer() {
	gormRepo, err := NewGormRepository(c.DBSourceURL)
	c.Require().NoError(err, "failed to initialize gorm database")
	c.Require().NotNil(gormRepo, "repository is nil")

	customer := &domain.Customer{}
	saveErr := gormRepo.Save(context.Background(), customer)
	c.NoError(saveErr, "failed to save customer")
}

func (c *CustomerRepoTestSuite) TestGetCustomerByName() {
	gormRepo, err := NewGormRepository(c.DBSourceURL)
	c.Require().NoError(err, "failed to initialize gorm database")
	c.Require().NotNil(gormRepo, "repository is nil")

	customer := &domain.Customer{
		Name:    "Alice",
		Status:  "active",
		Balance: 1500.0,
	}
	saveErr := gormRepo.Save(context.Background(), customer)
	c.Require().NoError(saveErr, "failed to save customer")

	retrievedCustomer, getErr := gormRepo.GetByName(context.Background(), customer.Name)
	c.NoError(getErr, "failed to get customer by name")

	c.Equal(customer.Name, retrievedCustomer.Name, "customer name does not match")
	c.Equal(customer.Status, retrievedCustomer.Status, "customer status does not match")
	c.Equal(customer.Balance, retrievedCustomer.Balance, "customer balance does not match")
}

func (c *CustomerRepoTestSuite) TestUpdateCustomer() {
	gormRepo, err := NewGormRepository(c.DBSourceURL)
	c.Require().NoError(err, "failed to initialize gorm database")
	c.Require().NotNil(gormRepo, "repository is nil")

	customer := &domain.Customer{
		Name:      "Jane",
		Status:    "inactive",
		Balance:   500.0,
		CreatedAt: time.Now(),
	}

	saveErr := gormRepo.Save(context.Background(), customer)
	c.Require().NoError(saveErr, "failed to save customer")

	// Update the customer status
	customer.Status = "active"
	updateErr := gormRepo.Update(context.Background(), customer)
	c.NoError(updateErr, "failed to update customer")

	updatedCustomer, getErr := gormRepo.GetByName(context.Background(), customer.Name)
	c.NoError(getErr, "failed to get updated customer")
	c.Equal("active", updatedCustomer.Status, "customer status was not updated correctly")
}

func TestCustomerRepoSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}
