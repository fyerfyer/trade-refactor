package repo

import (
	"context"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Payment struct {
	CustomerID uint64  `gorm:"not null;index"`
	OrderID    uint64  `gorm:"primaryKey;not null;index"`
	TotalPrice float32 `gorm:"not null"`
	Status     string  `gorm:"type:varchar(50);not null"`
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(dsn string) (*GormRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Payment{}); err != nil {
		return nil, err
	}

	return &GormRepository{db: db}, nil
}

func (r *GormRepository) Save(ctx context.Context, payment *domain.Payment) error {
	dbPayment := Payment{
		CustomerID: payment.CustomerID,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
		Status:     payment.Status,
	}

	return r.db.WithContext(ctx).Create(&dbPayment).Error
}
