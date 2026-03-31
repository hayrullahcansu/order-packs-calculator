// Package router defines HTTP handlers and route configuration for the API.
package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/cmd/api/requests"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/service"
)

// OrderPackHandler handles HTTP requests for order pack CRUD and calculation endpoints.
type OrderPackHandler struct {
	BaseHandler
	orderPackService *service.OrderPackService
}

// NewOrderPackHandler creates a handler with the given order pack service.
func NewOrderPackHandler(
	orderPackService *service.OrderPackService,
) (handler *OrderPackHandler) {
	handler = &OrderPackHandler{
		orderPackService: orderPackService,
	}
	return
}

// GetAvailableOrderPacks handles GET /v1/order_packs and returns all configured pack sizes.
func (handler *OrderPackHandler) GetAvailableOrderPacks(c *gin.Context) {
	orderPacks, fetchErr := handler.orderPackService.GetAvailableOrderPacks(c.Request.Context())
	if fetchErr != nil {
		handler.FailedBadRequest(c, nil, fetchErr)
		return
	}
	handler.OK(c, orderPacks)
}

// AddOrderPack handles POST /v1/order_packs and creates a new pack size.
func (handler *OrderPackHandler) AddOrderPack(c *gin.Context) {
	var req requests.AddOrderPackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.FailedBadRequest(c, nil, fmt.Errorf("invalid request body"))
		return
	}

	orderPack, err := handler.orderPackService.AddOrderPack(c.Request.Context(), req.Items)
	if err != nil {
		handler.FailedBadRequest(c, nil, err)
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Result: true, Data: orderPack})
}

// UpdateOrderPack handles PUT /v1/order_packs/:id and modifies an existing pack size.
func (handler *OrderPackHandler) UpdateOrderPack(c *gin.Context) {
	idParam := c.Param("id")
	id, parseErr := uuid.Parse(idParam)
	if parseErr != nil {
		handler.FailedBadRequest(c, nil, fmt.Errorf("invalid id"))
		return
	}

	var req requests.UpdateOrderPackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.FailedBadRequest(c, nil, fmt.Errorf("invalid request body"))
		return
	}

	orderPack, err := handler.orderPackService.UpdateOrderPack(c.Request.Context(), id, req.Items)
	if err != nil {
		handler.FailedBadRequest(c, nil, err)
		return
	}
	handler.OK(c, orderPack)
}

// RemoveOrderPack handles DELETE /v1/order_packs/:id and deletes a pack size.
func (handler *OrderPackHandler) RemoveOrderPack(c *gin.Context) {
	idParam := c.Param("id")
	id, parseErr := uuid.Parse(idParam)
	if parseErr != nil {
		handler.FailedBadRequest(c, nil, fmt.Errorf("invalid id"))
		return
	}

	if err := handler.orderPackService.RemoveOrderPack(c.Request.Context(), id); err != nil {
		handler.FailedBadRequest(c, nil, err)
		return
	}
	handler.OK(c, nil)
}

// SolveOrderPacks handles POST /v1/order_packs/solve and returns the optimal pack combination.
func (handler *OrderPackHandler) SolveOrderPacks(c *gin.Context) {
	var req requests.SolveOrderPacksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.FailedBadRequest(c, nil, fmt.Errorf("invalid request body"))
		return
	}

	orderPack, err := handler.orderPackService.SolveOrderPacks(c.Request.Context(), req.Order)
	if err != nil {
		handler.FailedBadRequest(c, nil, err)
		return
	}
	handler.OK(c, orderPack)
}
