package expenses

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repository) ListByUserId(ctx context.Context, userId int64) ([]Expense, error) {
	const query = `
		SELECT id, user_id, vendor, amount, category, description, date, created_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY date DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	expenses := make([]Expense, 0)
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(
			&expense.Id,
			&expense.UserId,
			&expense.Vendor,
			&expense.Amount,
			&expense.Category,
			&expense.Description,
			&expense.Date,
			&expense.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}

		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expenses: %w", err)
	}

	return expenses, nil
}

func (r *Repository) GetTotalByUserId(ctx context.Context, userId int64) (float64, error) {
	const query = `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = $1
	`

	var total float64
	if err := r.db.QueryRowContext(ctx, query, userId).Scan(&total); err != nil {
		return 0, fmt.Errorf("get total expenses: %w", err)
	}

	return total, nil
}

func (r *Repository) GetById(ctx context.Context, expenseId, userId int64) (Expense, error) {
	const query = `
		SELECT id, user_id, vendor, amount, category, description, date, created_at
		FROM expenses
		WHERE id = $1 AND user_id = $2
	`

	var expense Expense
	err := r.db.QueryRowContext(ctx, query, expenseId, userId).Scan(
		&expense.Id,
		&expense.UserId,
		&expense.Vendor,
		&expense.Amount,
		&expense.Category,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Expense{}, ErrExpenseNotFound
		}
		return Expense{}, fmt.Errorf("get expense by id: %w", err)
	}

	return expense, nil
}

func (r *Repository) Update(ctx context.Context, expense Expense) error {
	const query = `
		UPDATE expenses
		SET vendor = $1, amount = $2, category = $3, description = $4, date = $5
		WHERE id = $6 AND user_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		expense.Vendor,
		expense.Amount,
		expense.Category,
		expense.Description,
		expense.Date,
		expense.Id,
		expense.UserId,
	)
	if err != nil {
		return fmt.Errorf("update expense: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get updated rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrExpenseNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, expenseId, userId int64) error {
	const query = `
		DELETE FROM expenses
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, expenseId, userId)
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrExpenseNotFound
	}

	return nil
}
