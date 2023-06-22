package controller

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required,alpha"`
	Currency string `json:"currency" binding:"required,oneof=EUR BGN USD"`
}

type ListAccountQuery struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=1"`
}
