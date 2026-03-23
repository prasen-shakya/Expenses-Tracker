package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	params, err := GetRequestParams[expenses.ExpenseParams](w, r)

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
		Date:        parseExpenseDate(params.Date),
	}

	err = e.expenseService.LogExpense(r.Context(), expense)

	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to log expense"})
		return
	}

	respond.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Logged expense successfully"})
}

func (e *ExpenseController) ListExpenses(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	expenseList, err := e.expenseService.ListExpenses(r.Context(), int64(user.Id))
	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to list expenses"})
		return
	}

	respond.WriteJSON(w, http.StatusOK, expenseList)
}

func (e *ExpenseController) GetTotalExpenses(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	total, err := e.expenseService.GetTotalExpenses(r.Context(), int64(user.Id))
	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to get total expenses"})
		return
	}

	respond.WriteJSON(w, http.StatusOK, map[string]float64{"total": total})
}

func (e *ExpenseController) GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	expenseId, err := getExpenseID(r)
	if err != nil {
		respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid expense id"})
		return
	}

	expense, err := e.expenseService.GetExpenseByID(r.Context(), expenseId, int64(user.Id))
	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if errors.Is(err, expenses.ErrExpenseNotFound) {
			respond.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Expense not found"})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to get expense"})
		return
	}

	respond.WriteJSON(w, http.StatusOK, expense)
}

func (e *ExpenseController) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	params, err := GetRequestParams[expenses.ExpenseParams](w, r)
	if err != nil {
		return
	}

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	expenseId, err := getExpenseID(r)
	if err != nil {
		respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid expense id"})
		return
	}

	expense := expenses.Expense{
		Id:          expenseId,
		UserId:      int64(user.Id),
		Vendor:      params.Vendor,
		Amount:      params.Amount,
		Category:    params.Category,
		Description: params.Description,
		Date:        parseExpenseDate(params.Date),
	}

	err = e.expenseService.UpdateExpense(r.Context(), expense)
	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if errors.Is(err, expenses.ErrExpenseNotFound) {
			respond.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Expense not found"})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to update expense"})
		return
	}

	respond.WriteJSON(w, http.StatusOK, map[string]string{"message": "Updated expense successfully"})
}

func (e *ExpenseController) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		return
	}

	expenseId, err := getExpenseID(r)
	if err != nil {
		respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid expense id"})
		return
	}

	err = e.expenseService.DeleteExpense(r.Context(), expenseId, int64(user.Id))
	if err != nil {
		if isExpenseValidationError(err) {
			respond.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if errors.Is(err, expenses.ErrExpenseNotFound) {
			respond.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Expense not found"})
			return
		}

		respond.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to delete expense"})
		return
	}

	respond.WriteJSON(w, http.StatusOK, map[string]string{"message": "Deleted expense successfully"})
}

func getExpenseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func parseExpenseDate(value string) time.Time {
	if value == "" {
		return time.Now().UTC()
	}

	parsedTime, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return parsedTime
	}

	return time.Now().UTC()
}

func isExpenseValidationError(err error) bool {
	return errors.Is(err, expenses.ErrInvalidExpenseUserID) ||
		errors.Is(err, expenses.ErrInvalidExpenseID) ||
		errors.Is(err, expenses.ErrInvalidExpenseVendor) ||
		errors.Is(err, expenses.ErrInvalidExpenseCategory) ||
		errors.Is(err, expenses.ErrInvalidExpenseAmount)
}
