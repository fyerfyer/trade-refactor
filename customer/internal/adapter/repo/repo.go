package repo

import (
	"context"

	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name    string  `gorm:"uniqueIndex;type:varchar(100);not null"`
	Status  string  `gorm:"type:varchar(50);not null"`
	Balance float32 `gorm:"not null"`
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(dsn string) (*GormRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Customer{}); err != nil {
		return nil, err
	}
	return &GormRepository{db: db}, nil
}

func (r *GormRepository) Save(ctx context.Context, customer *domain.Customer) error {
	dbCustomer := convertDomainCustomerIntoDB(*customer)
	err := r.db.WithContext(ctx).Create(&dbCustomer).Error
	if err == nil {
		customer.ID = uint64(dbCustomer.ID)
	}

	return err
}

func (r *GormRepository) Update(ctx context.Context, customer *domain.Customer) error {
	dbCustomer := convertDomainCustomerIntoDB(*customer)
	return r.db.WithContext(ctx).Save(&dbCustomer).Error
}

// func (r *GormRepository) GetByID(ctx context.Context, customerID uint64) (*domain.Customer, error) {
// 	var customer Customer
// 	err := r.db.WithContext(ctx).Where("id = ?", customerID).First(&customer).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	domainCustomer := ConvertDBIntoDomainCustomer(customer)
// 	return &domainCustomer, nil
// }

func (r *GormRepository) GetByName(ctx context.Context, customerName string) (*domain.Customer, error) {
	var customer Customer
	err := r.db.WithContext(ctx).Where("name = ?", customerName).First(&customer).Error
	if err != nil {
		return nil, err
	}

	domainCustomer := convertDBIntoDomainCustomer(customer)
	return &domainCustomer, nil
}

func convertDomainCustomerIntoDB(customer domain.Customer) Customer {
	return Customer{
		Model: gorm.Model{
			ID:        uint(customer.ID),
			CreatedAt: customer.CreatedAt,
		},
		Name:    customer.Name,
		Status:  customer.Status,
		Balance: customer.Balance,
	}
}

func convertDBIntoDomainCustomer(customer Customer) domain.Customer {
	return domain.Customer{
		ID:        uint64(customer.ID),
		Name:      customer.Name,
		Status:    customer.Status,
		Balance:   customer.Balance,
		CreatedAt: customer.CreatedAt,
	}
}
