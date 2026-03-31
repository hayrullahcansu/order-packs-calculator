package model_test

import (
	"testing"

	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidate_NilOrderPack(t *testing.T) {
	var op *model.OrderPack
	err := op.Validate()
	assert.Error(t, err)
	assert.Equal(t, "order pack cannot be nil", err.Error())
}

func TestValidate_ZeroItems(t *testing.T) {
	op := &model.OrderPack{Items: 0}
	err := op.Validate()
	assert.Error(t, err)
	assert.Equal(t, "invalid items", err.Error())
}

func TestValidate_NegativeItems(t *testing.T) {
	op := &model.OrderPack{Items: -1}
	err := op.Validate()
	assert.Error(t, err)
	assert.Equal(t, "invalid items", err.Error())
}

func TestValidate_ValidItems(t *testing.T) {
	op := &model.OrderPack{Items: 250}
	err := op.Validate()
	assert.NoError(t, err)
}

func TestValidate_MinimumValidItems(t *testing.T) {
	op := &model.OrderPack{Items: 1}
	err := op.Validate()
	assert.NoError(t, err)
}
