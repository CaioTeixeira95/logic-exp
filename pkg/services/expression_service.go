package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/CaioTeixeira95/logic-exp/pkg/repositories"
	"github.com/CaioTeixeira95/logic-exp/pkg/utils"
)

var ErrInvalidExpression = errors.New("invalid expression")

type ExpressionService interface {
	CreateExpression(ctx context.Context, exp *repositories.Expression) (*repositories.Expression, error)
	ListExpressions(ctx context.Context) ([]repositories.Expression, error)
	UpdateExpression(ctx context.Context, exp *repositories.Expression) (*repositories.Expression, error)
	EvaluateExpression(ctx context.Context, ID int64, parameters map[string]int) (bool, error)
}

type expressionService struct {
	expressionRepository repositories.ExpressionRepository
}

func (es *expressionService) CreateExpression(ctx context.Context, exp *repositories.Expression) (*repositories.Expression, error) {
	if err := validateExpression(exp.Value); err != nil {
		return nil, err
	}

	exp, err := es.expressionRepository.CreateExpression(ctx, exp)
	if err != nil {
		return nil, fmt.Errorf("error creating expression: %w", err)
	}

	return exp, nil
}

func (es *expressionService) ListExpressions(ctx context.Context) ([]repositories.Expression, error) {
	exps, err := es.expressionRepository.GetAllExpressions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all expressions: %w", err)
	}

	return exps, nil
}

func (es *expressionService) UpdateExpression(ctx context.Context, exp *repositories.Expression) (*repositories.Expression, error) {
	if exp.ID == 0 {
		return nil, fmt.Errorf("invalid expression ID provided")
	}

	if err := validateExpression(exp.Value); err != nil {
		return nil, err
	}

	updatedExp, err := es.expressionRepository.UpdateExpression(ctx, exp)
	if err == repositories.ErrNoRowsAffected {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("error updating expression ID %d: %w", exp.ID, err)
	}

	return updatedExp, nil
}

func validateExpression(expression string) error {
	if expression == "" {
		return fmt.Errorf("value can't be empty")
	}

	isValid := utils.IsLogicalExpressionValid(expression)
	if !isValid {
		return ErrInvalidExpression
	}

	return nil
}

func (es *expressionService) EvaluateExpression(ctx context.Context, ID int64, parameters map[string]int) (bool, error) {
	exp, err := es.expressionRepository.GetExpressionByID(ctx, ID)
	if err == repositories.ErrExpressionNotFound {
		return false, err
	}
	if err != nil {
		return false, fmt.Errorf("error getting expression ID %d: %w", ID, err)
	}

	// Validate if all expected parameters were provided
	expExpectedParameters := utils.GetLogicalExpressionParameters(exp.Value)
	for key := range expExpectedParameters {
		if _, ok := parameters[key]; !ok {
			return false, fmt.Errorf("missing parameter %q for the logical expression %q", key, exp.Value)
		}
	}

	res, err := utils.EvaluateLogicalExpression(exp.Value, parameters)
	if err != nil {
		return false, fmt.Errorf("error evaluating expression %q: %w", exp.Value, err)
	}

	return res, nil
}

type ExpressionServiceOption func(es *expressionService)

func NewExpressionService(options ...ExpressionServiceOption) ExpressionService {
	es := &expressionService{}

	for _, option := range options {
		option(es)
	}

	return es
}

var _ ExpressionService = (*expressionService)(nil)

func WithExpressionRepositoryOption(er repositories.ExpressionRepository) ExpressionServiceOption {
	return func(es *expressionService) {
		es.expressionRepository = er
	}
}
