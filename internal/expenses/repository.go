package expenses

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, expense Expense) error {
	const query = `
		INSERT INTO expenses (user_id, vendor, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		expense.UserId,
		expense.Vendor,
		expense.Amount,
		expense.Category,
		expense.Description,
		expense.Date,
	)
	if err != nil {
		return fmt.Errorf("create expense: %w", err)
	}

	return nil
}
