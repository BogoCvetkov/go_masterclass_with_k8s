package api

import "time"

// Account

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required,alpha"`
	Currency string `json:"currency" binding:"required,validcurrency"`
}

type ListAccountQuery struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=1"`
}

// Transfer

type CreateTransferParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,number"`
	ToAccountID   int64  `json:"to_account_id"  binding:"required,number"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,validcurrency"`
}

// User
type CreateUserParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required,alpha"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserReponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
}

type LoginUserParams struct {
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}
