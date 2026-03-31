package requests

type AddOrderPackRequest struct {
	Items int `json:"items" binding:"required"`
}

type UpdateOrderPackRequest struct {
	Items int `json:"items" binding:"required"`
}

type SolveOrderPacksRequest struct {
	Order int `json:"order" binding:"required"`
}
