package expenses

import (
	"context"
	"strings"
)

type Service struct {
	expensesRepo *Repository
}

func NewService(expensesRepo *Repository) *Service {
	return &Service{expensesRepo: expensesRepo}
}

func (s *Service) LogExpense(ctx context.Context, expense Expense) error {
	if err := validateExpense(expense); err != nil {
		return err
	}

	expense.Vendor = strings.TrimSpace(expense.Vendor)
	expense.Category = strings.TrimSpace(expense.Category)

	return s.expensesRepo.Create(ctx, expense)
}

func (s *Service) ListExpenses(ctx context.Context, userId int64) ([]Expense, error) {
	if userId <= 0 {
		return nil, ErrInvalidExpenseUserID
	}

	return s.expensesRepo.ListByUserId(ctx, userId)
}

func (s *Service) GetTotalExpenses(ctx context.Context, userId int64) (float64, error) {
	if userId <= 0 {
		return 0, ErrInvalidExpenseUserID
	}

	return s.expensesRepo.GetTotalByUserId(ctx, userId)
}

func (s *Service) GetExpenseByID(ctx context.Context, expenseId, userId int64) (Expense, error) {
	if expenseId <= 0 {
		return Expense{}, ErrInvalidExpenseID
	}
	if userId <= 0 {
		return Expense{}, ErrInvalidExpenseUserID
	}

	return s.expensesRepo.GetById(ctx, expenseId, userId)
}

func (s *Service) UpdateExpense(ctx context.Context, expense Expense) error {
	if expense.Id <= 0 {
		return ErrInvalidExpenseID
	}
	if err := validateExpense(expense); err != nil {
		return err
	}

	expense.Vendor = strings.TrimSpace(expense.Vendor)
	expense.Category = strings.TrimSpace(expense.Category)

	return s.expensesRepo.Update(ctx, expense)
}

func (s *Service) DeleteExpense(ctx context.Context, expenseId, userId int64) error {
	if expenseId <= 0 {
		return ErrInvalidExpenseID
	}
	if userId <= 0 {
		return ErrInvalidExpenseUserID
	}

	return s.expensesRepo.Delete(ctx, expenseId, userId)
}

func validateExpense(expense Expense) error {
	if expense.UserId <= 0 {
		return ErrInvalidExpenseUserID
	}
	if strings.TrimSpace(expense.Vendor) == "" {
		return ErrInvalidExpenseVendor
	}
	if strings.TrimSpace(expense.Category) == "" {
		return ErrInvalidExpenseCategory
	}
	if expense.Amount <= 0 {
		return ErrInvalidExpenseAmount
	}

	return nil
}
