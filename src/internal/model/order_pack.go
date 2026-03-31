package model

import "fmt"

// OrderPack represents an available pack size that can be used to fulfill orders.
type OrderPack struct {
	TimeAwareEntity
	Items int `gorm:"uniqueIndex" json:"items" example:"250"`
}

func (op *OrderPack) Validate() error {
	if op == nil {
		return fmt.Errorf("order pack cannot be nil")
	}

	if op.Items <= 0 {
		return fmt.Errorf("invalid items")
	}
	return nil
}
