package expenses

import (
	"context"
)

type Service struct {
	expensesRepo *Repository
}

func NewService(expensesRepo *Repository) *Service {
	return &Service{expensesRepo: expensesRepo}
}

func (s *Service) LogExpense(ctx context.Context, params Expense) error {
	err := s.expensesRepo.Create(ctx, params)

	if err != nil {
		return err
	}
	return nil
}
