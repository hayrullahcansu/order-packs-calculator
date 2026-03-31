package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
)

type OrderPackRepository interface {
	FetchAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error)
	AddOrderPack(ctx context.Context, orderPack *model.OrderPack) error
	GetOrderPackByID(ctx context.Context, id uuid.UUID) (*model.OrderPack, error)
	UpdateOrderPack(ctx context.Context, updatedOrderPack *model.OrderPack) error
	RemoveOrderPack(ctx context.Context, id uuid.UUID) error
}

type OrderPackService struct {
	orderPackRepository OrderPackRepository
	orderPackCalculator OrderPackCalculator
}

func NewOrderPackService(
	orderPackRepository OrderPackRepository,
) *OrderPackService {
	manager := &OrderPackService{
		orderPackRepository: orderPackRepository,
		orderPackCalculator: NewOrderPackCalculator(),
	}
	return manager
}

func (s *OrderPackService) GetAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.orderPackRepository.FetchAvailableOrderPacks(ctx)
}

func (s *OrderPackService) AddOrderPack(ctx context.Context, items int) (*model.OrderPack, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	orderPack := &model.OrderPack{Items: items}

	if validated := orderPack.Validate(); validated != nil {
		return nil, validated
	}

	insertErr := s.orderPackRepository.AddOrderPack(ctx, orderPack)
	if insertErr != nil {
		return nil, fmt.Errorf("duplicated order pack")
	}
	return orderPack, nil
}

func (s *OrderPackService) UpdateOrderPack(ctx context.Context, id uuid.UUID, items int) (*model.OrderPack, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	orderPack, fetchErr := s.orderPackRepository.GetOrderPackByID(ctx, id)
	if fetchErr != nil {
		return nil, fmt.Errorf("fetching error")
	}

	if orderPack == nil {
		return nil, fmt.Errorf("order pack not found")
	}
	orderPack.Items = items
	if validated := orderPack.Validate(); validated != nil {
		return nil, validated
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	updateErr := s.orderPackRepository.UpdateOrderPack(ctx, orderPack)
	if updateErr != nil {
		return nil, fmt.Errorf("duplicated order pack")
	}
	return orderPack, nil
}

func (s *OrderPackService) RemoveOrderPack(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	updateErr := s.orderPackRepository.RemoveOrderPack(ctx, id)
	if updateErr != nil {
		return fmt.Errorf("order couldn't removed or exists")
	}
	return nil
}

func (s *OrderPackService) SolveOrderPacks(ctx context.Context, order int) (map[int]int, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	orderPacks, fetchErr := s.orderPackRepository.FetchAvailableOrderPacks(ctx)
	if fetchErr != nil {
		return nil, fetchErr
	}
	if len(orderPacks) == 0 {
		return nil, fmt.Errorf("no available order packs")
	}

	packs := make([]int, len(orderPacks))
	for i, v := range orderPacks {
		packs[i] = v.Items
	}
	result := s.orderPackCalculator.SolvePacks(packs, order)
	if result == nil {
		return nil, fmt.Errorf("no available order packs")
	}
	return result, nil
}
