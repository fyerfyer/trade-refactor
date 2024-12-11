package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PaymentRepoTestSuite struct {
	suite.Suite
	DBSourceURL string
}

func (p *PaymentRepoTestSuite) SetupSuite() {
	ctx := context.Background()
	port := nat.Port("3306/tcp")

	// Define database URL function
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf(
			"root:110119abc@tcp(%s:%s)/orders?charset=utf8mb4&parseTime=True&loc=Local",
			host, port.Port(),
		)
	}

	// Define container request
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{string(port)},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "110119abc",
			"MYSQL_DATABASE":      "orders",
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
	p.Require().NoError(err, "failed to start MySQL container")

	// Get host, port from the container
	host, err := mysqlContainer.Host(ctx)
	p.Require().NoError(err, "failed to get container host")

	mappedPort, err := mysqlContainer.MappedPort(ctx, port)
	p.Require().NoError(err, "failed to get container mapped port")

	p.DBSourceURL = dbURL(host, mappedPort)
}

func (p *PaymentRepoTestSuite) Test_Should_Save_Payment() {
	gormRepo, err := NewGormRepository(p.DBSourceURL)
	p.Require().NoError(err, "failed to initialize gorm database")
	p.Require().NotNil(gormRepo, "repository is nil")

	payment := &domain.Payment{}
	saveErr := gormRepo.Save(context.Background(), payment)
	p.NoError(saveErr, "failed to save order")
}

func TestPaymentRepoSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepoTestSuite))
}
