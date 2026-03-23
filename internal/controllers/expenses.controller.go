package controllers

import (
	"net/http"

	"github.com/prasen-shakya/todo/internal/auth"
	"github.com/prasen-shakya/todo/internal/expenses"
	"github.com/prasen-shakya/todo/internal/respond"
)

type ExpenseController struct {
	expenseService *expenses.Service
}

func NewExpenseController(expenseService *expenses.Service) *ExpenseController {
	return &ExpenseController{expenseService: expenseService}
}

func (e *ExpenseController) LogExpense(w http.ResponseWriter, r *http.Request) {
	params, err := GetRequestParams[expenses.LogExpenseParams](w, r)

	if err != nil {
		return
	}

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	expense := expenses.Expense{
		UserId:      int64(user.Id),
		Vendor:      params.Vendor,
		Amount:      params.Amount,
		Category:    params.Category,
		Description: params.Description,
		Date:        params.Date,
	}

	err = e.expenseService.LogExpense(r.Context(), expense)

	if err != nil {
		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to log expense"})
		return
	}

	respond.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Logged expense successfully"})
}
