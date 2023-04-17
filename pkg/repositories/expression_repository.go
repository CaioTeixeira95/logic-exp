package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrExpressionNotFound = errors.New("expression not found")
	ErrNoRowsAffected     = errors.New("no rows affected")
)

type Expression struct {
	ID    int64
	Value string
}

type ExpressionRepository interface {
	GetAllExpressions(ctx context.Context) ([]Expression, error)
	GetExpressionByID(ctx context.Context, ID int64) (*Expression, error)
	UpdateExpression(ctx context.Context, exp *Expression) (*Expression, error)
	CreateExpression(ctx context.Context, exp *Expression) (*Expression, error)
}

func (r *DefaultRepository) GetAllExpressions(ctx context.Context) ([]Expression, error) {
	const query = `
		SELECT
			id,
			expression
		FROM
			expressions
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying expressions: %w", err)
	}
	defer rows.Close()

	exps := []Expression{}
	for rows.Next() {
		var ID int64
		var value string
		if err := rows.Scan(&ID, &value); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		exps = append(exps, Expression{
			ID:    ID,
			Value: value,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %w", err)
	}

	return exps, nil
}

func (r *DefaultRepository) GetExpressionByID(ctx context.Context, ID int64) (*Expression, error) {
	const query = `
		SELECT
			expression
		FROM
			expressions
		WHERE
			id = $1
	`

	var value string
	err := r.db.QueryRowContext(ctx, query, ID).Scan(&value)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying expression ID %d: %w", ID, err)
	}
	if err != nil {
		return nil, ErrExpressionNotFound
	}

	return &Expression{
		ID:    ID,
		Value: value,
	}, nil
}

func (r *DefaultRepository) CreateExpression(ctx context.Context, exp *Expression) (*Expression, error) {
	const query = `
		INSERT INTO expressions
			(expression)
		VALUES
			($1)
		RETURNING id
	`
	var ID int64
	err := r.db.QueryRowContext(ctx, query, exp.Value).Scan(&ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting new expression: %w", err)
	}

	exp.ID = ID

	return exp, nil
}

func (r *DefaultRepository) UpdateExpression(ctx context.Context, exp *Expression) (*Expression, error) {
	const query = `
		UPDATE
			expressions
		SET
			expression = $2
		WHERE
			id = $1
	`
	result, err := r.db.ExecContext(ctx, query, exp.ID, exp.Value)
	if err != nil {
		return nil, fmt.Errorf("error updating expression ID %d: %w", exp.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting number of rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNoRowsAffected
	}

	return exp, nil
}
