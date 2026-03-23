package routes

import (
	"net/http"

	"github.com/prasen-shakya/todo/internal/auth"
	"github.com/prasen-shakya/todo/internal/controllers"
	"github.com/prasen-shakya/todo/internal/users"
)

func RegisterExpenseRoutes(
	mux *http.ServeMux,
	expenseController *controllers.ExpenseController,
	authService *auth.Service,
	usersRepo *users.Repository,
) {

	requireAuth := auth.RequireAuth(authService, usersRepo)

	mux.Handle("POST /expenses", requireAuth(http.HandlerFunc(expenseController.LogExpense)))
	mux.Handle("GET /expenses", requireAuth(http.HandlerFunc(expenseController.ListExpenses)))
	mux.Handle("GET /expenses/total", requireAuth(http.HandlerFunc(expenseController.GetTotalExpenses)))
	mux.Handle("GET /expenses/{id}", requireAuth(http.HandlerFunc(expenseController.GetExpenseByID)))
	mux.Handle("PUT /expenses/{id}", requireAuth(http.HandlerFunc(expenseController.UpdateExpense)))
	mux.Handle("DELETE /expenses/{id}", requireAuth(http.HandlerFunc(expenseController.DeleteExpense)))
}
