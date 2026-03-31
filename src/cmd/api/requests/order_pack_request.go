// Package requests defines the JSON request payloads for the order pack API.
package requests

// AddOrderPackRequest is the payload for creating a new pack size.
type AddOrderPackRequest struct {
	Items int `json:"items" binding:"required"`
}

// UpdateOrderPackRequest is the payload for modifying an existing pack size.
type UpdateOrderPackRequest struct {
	Items int `json:"items" binding:"required"`
}

// SolveOrderPacksRequest is the payload for calculating optimal packs for a given order quantity.
type SolveOrderPacksRequest struct {
	Order int `json:"order" binding:"required"`
}
