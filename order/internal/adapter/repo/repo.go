package repo

import (
	"context"
	"fmt"

	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID uint64      `gorm:"not null;index" json:"customer_id"`
	Status     string      `gorm:"type:varchar(50);not null" json:"status"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items"`
}

type OrderItem struct {
	ProductCode string  `gorm:"primaryKey;type:varchar(100);not null" json:"product_code"`
	UnitPrice   float32 `gorm:"not null" json:"unit_price"`
	Quantity    int32   `gorm:"not null" json:"quantity"`
	OrderID     uint64  `gorm:"primaryKey;not null;index" json:"order_id"`
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(dsn string) (*GormRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// return nil, err
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&Order{}); err != nil {
		// return nil, err
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&OrderItem{}); err != nil {
		// return nil, err
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &GormRepository{db: db}, nil
}

func (r *GormRepository) Save(ctx context.Context, order *domain.Order) error {
	// convert domain object into database object
	orderModel := convertDomainOrderIntoDB(*order)

	res := r.db.WithContext(ctx).Create(&orderModel)
	if res.Error == nil {
		order.ID = uint64(orderModel.ID)
	}
	return res.Error
}

func (r *GormRepository) Update(ctx context.Context, order *domain.Order) error {
	// convert domain object into database object
	orderModel := convertDomainOrderIntoDB(*order)
	res := r.db.WithContext(ctx).Save(&orderModel)
	return res.Error
}

func (r *GormRepository) Delete(ctx context.Context, orderID uint64) error {
	return r.db.WithContext(ctx).Delete(&Order{}, orderID).Error
}

func (r *GormRepository) GetUnpaidOrders(ctx context.Context, customerID uint64) ([]domain.Order, error) {
	var orders []Order
	err := r.db.WithContext(ctx).Where("customer_id = ? AND status = ?", customerID, "unpaid").
		Find(&orders).
		Error

	if err != nil {
		return nil, err
	}

	// convert database model into domain model
	return convertDBIntoDomainOrders(orders), nil
}

func (r *GormRepository) GetUnpaidOrder(ctx context.Context, orderID uint64) (*domain.Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Where("id = ? AND status = ?", orderID, "unpaid").
		First(&o).
		Error

	if err != nil {
		return nil, err
	}

	domainOrder := convertDBIntoDomainOrder(o)
	return &domainOrder, nil
}

func convertDomainOrderIntoDB(order domain.Order) Order {
	var orderItems []OrderItem
	for _, orderItem := range order.Items {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	return Order{
		Model: gorm.Model{
			ID:        uint(order.ID),
			CreatedAt: order.CreatedAt,
		},
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
	}
}

func convertDomainOrdersIntoDB(orders []domain.Order) []Order {
	var dbOrders []Order
	for _, order := range orders {
		dbOrders = append(dbOrders, convertDomainOrderIntoDB(order))
	}

	return dbOrders
}

func convertDBIntoDomainOrder(order Order) domain.Order {
	var items []domain.OrderItem
	for _, orderItem := range order.OrderItems {
		items = append(items, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	return domain.Order{
		ID:         uint64(order.ID),
		CustomerID: order.CustomerID,
		Items:      items,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}
}

func convertDBIntoDomainOrders(orders []Order) []domain.Order {
	var domainOrders []domain.Order
	for _, order := range orders {
		domainOrders = append(domainOrders, convertDBIntoDomainOrder(order))
	}

	return domainOrders
}
