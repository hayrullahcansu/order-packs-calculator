// Package repository provides the data access layer for order pack persistence using GORM.
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
	"gorm.io/gorm"
)

// OrderPackRepository implements the service.OrderPackRepository interface
// using GORM with SQLite as the underlying database.
type OrderPackRepository struct {
	db *gorm.DB
}

// NewOrderPackRepository creates a new repository backed by the given GORM database connection.
func NewOrderPackRepository(db *gorm.DB) *OrderPackRepository {
	return &OrderPackRepository{db: db}
}

// FetchAvailableOrderPacks retrieves all order packs from the database.
func (r *OrderPackRepository) FetchAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error) {
	var packs []*model.OrderPack
	result := r.db.WithContext(ctx).Find(&packs)
	return packs, result.Error
}

// AddOrderPack inserts a new order pack. Fails if the items value already exists (unique index).
func (r *OrderPackRepository) AddOrderPack(ctx context.Context, orderPack *model.OrderPack) error {
	return r.db.WithContext(ctx).Create(orderPack).Error
}

// GetOrderPackByID retrieves a single order pack by its UUID.
func (r *OrderPackRepository) GetOrderPackByID(ctx context.Context, id uuid.UUID) (*model.OrderPack, error) {
	var pack model.OrderPack
	result := r.db.WithContext(ctx).First(&pack, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pack, nil
}

// UpdateOrderPack saves all fields of the given order pack back to the database.
func (r *OrderPackRepository) UpdateOrderPack(ctx context.Context, updatedOrderPack *model.OrderPack) error {
	return r.db.WithContext(ctx).Save(updatedOrderPack).Error
}

// RemoveOrderPack deletes an order pack by its UUID.
func (r *OrderPackRepository) RemoveOrderPack(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.OrderPack{}, "id = ?", id).Error
}
