package expenses

import (
	"errors"
	"time"
)

var ErrExpenseNotFound = errors.New("expense not found")
var ErrInvalidExpenseUserID = errors.New("user id is required")
var ErrInvalidExpenseID = errors.New("expense id is required")
var ErrInvalidExpenseVendor = errors.New("vendor is required")
var ErrInvalidExpenseCategory = errors.New("category is required")
var ErrInvalidExpenseAmount = errors.New("amount must be greater than 0")

type ExpenseParams struct {
	Vendor      string  `json:"vendor"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type Expense struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Vendor      string    `json:"vendor"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}
