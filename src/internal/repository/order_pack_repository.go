package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
	"gorm.io/gorm"
)

type OrderPackRepository struct {
	db *gorm.DB
}

func NewOrderPackRepository(db *gorm.DB) *OrderPackRepository {
	return &OrderPackRepository{db: db}
}

func (r *OrderPackRepository) FetchAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error) {
	var packs []*model.OrderPack
	result := r.db.WithContext(ctx).Find(&packs)
	return packs, result.Error
}

func (r *OrderPackRepository) AddOrderPack(ctx context.Context, orderPack *model.OrderPack) error {
	return r.db.WithContext(ctx).Create(orderPack).Error
}

func (r *OrderPackRepository) GetOrderPackByID(ctx context.Context, id uuid.UUID) (*model.OrderPack, error) {
	var pack model.OrderPack
	result := r.db.WithContext(ctx).First(&pack, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pack, nil
}

func (r *OrderPackRepository) UpdateOrderPack(ctx context.Context, updatedOrderPack *model.OrderPack) error {
	return r.db.WithContext(ctx).Save(updatedOrderPack).Error
}

func (r *OrderPackRepository) RemoveOrderPack(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.OrderPack{}, "id = ?", id).Error
}
