// Package service contains the business logic for order pack management and calculation.
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
)

// OrderPackRepository defines the data access contract for order pack persistence.
type OrderPackRepository interface {
	FetchAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error)
	AddOrderPack(ctx context.Context, orderPack *model.OrderPack) error
	GetOrderPackByID(ctx context.Context, id uuid.UUID) (*model.OrderPack, error)
	UpdateOrderPack(ctx context.Context, updatedOrderPack *model.OrderPack) error
	RemoveOrderPack(ctx context.Context, id uuid.UUID) error
}

// OrderPackCalculator defines the interface for solving order pack optimization problems.
type OrderPackCalculator interface {
	// SolvePacks finds the minimum number of packs that meet or exceed the given order quantity.
	// It returns a map of pack size to count. Example: {500: 1, 250: 1} for order=501.
	SolvePacks(packs []int, order int) map[int]int
}

// OrderPackService provides CRUD operations for pack sizes and delegates
// order calculations to the OrderPackCalculator.
type OrderPackService struct {
	orderPackRepository OrderPackRepository
	orderPackCalculator OrderPackCalculator
}

// NewOrderPackService creates a new service with the given repository
// and an internally initialized calculator.
func NewOrderPackService(
	orderPackRepository OrderPackRepository,
) *OrderPackService {
	manager := &OrderPackService{
		orderPackRepository: orderPackRepository,
		orderPackCalculator: NewOrderPackCalculator(),
	}
	return manager
}

// GetAvailableOrderPacks returns all configured pack sizes from the database.
func (s *OrderPackService) GetAvailableOrderPacks(ctx context.Context) ([]*model.OrderPack, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.orderPackRepository.FetchAvailableOrderPacks(ctx)
}

// AddOrderPack validates and persists a new pack size. Returns an error if the size is duplicate.
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

// UpdateOrderPack modifies an existing pack size identified by its UUID.
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

// RemoveOrderPack deletes a pack size by its UUID.
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

// maxOrderLimit caps the order size to prevent excessive memory allocation
// in the dynamic programming algorithm (~16 bytes per unit).
const maxOrderLimit = 10_000_000

// SolveOrderPacks calculates the optimal pack combination for a given order quantity.
// It returns a map of pack size to count, minimizing both total items and number of packs.
func (s *OrderPackService) SolveOrderPacks(ctx context.Context, order int) (map[int]int, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if order <= 0 {
		return nil, fmt.Errorf("order must be greater than 0")
	}
	if order > maxOrderLimit {
		return nil, fmt.Errorf("order must not exceed %d", maxOrderLimit)
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
