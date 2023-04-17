package services

import (
	"context"
	"errors"
	"testing"

	"github.com/CaioTeixeira95/logic-exp/pkg/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpressionService_CreateExpression(t *testing.T) {
	expressionRepositoryMock := &repositories.ExpressionRepositoryMock{}
	expressionService := NewExpressionService(WithExpressionRepositoryOption(expressionRepositoryMock))

	ctx := context.Background()

	t.Run("returns error when the expressions is invalid", func(t *testing.T) {
		exp, err := expressionService.CreateExpression(ctx, &repositories.Expression{
			Value: "",
		})

		assert.EqualError(t, err, "value can't be empty")
		assert.Nil(t, exp)

		exp, err = expressionService.CreateExpression(ctx, &repositories.Expression{
			Value: "AND",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)

		exp, err = expressionService.CreateExpression(ctx, &repositories.Expression{
			Value: "x AND",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)

		exp, err = expressionService.CreateExpression(ctx, &repositories.Expression{
			Value: "x AND b OR",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		exp := &repositories.Expression{
			Value: "x AND y",
		}

		expressionRepositoryMock.
			On("CreateExpression", ctx, exp).
			Return(nil, errors.New("unexpected error")).
			Once()

		exp, err := expressionService.CreateExpression(ctx, exp)

		assert.EqualError(t, err, "error creating expression: unexpected error")
		assert.Nil(t, exp)
	})

	t.Run("creates an expression correctly", func(t *testing.T) {
		exp := &repositories.Expression{
			Value: "x AND y",
		}

		expectedExp := &repositories.Expression{
			ID:    1,
			Value: "x AND y",
		}

		expressionRepositoryMock.
			On("CreateExpression", ctx, exp).
			Return(expectedExp, nil).
			Once()

		exp, err := expressionService.CreateExpression(ctx, exp)
		require.NoError(t, err)

		assert.Equal(t, expectedExp, exp)
	})

	expressionRepositoryMock.AssertExpectations(t)
}

func TestExpressionService_ListExpressions(t *testing.T) {
	expressionRepositoryMock := &repositories.ExpressionRepositoryMock{}
	expressionService := NewExpressionService(WithExpressionRepositoryOption(expressionRepositoryMock))

	ctx := context.Background()

	t.Run("returns error when repository fails", func(t *testing.T) {
		expressionRepositoryMock.
			On("GetAllExpressions", ctx).
			Return(nil, errors.New("unexpected error")).
			Once()

		exps, err := expressionService.ListExpressions(ctx)

		assert.EqualError(t, err, "error getting all expressions: unexpected error")
		assert.Nil(t, exps)
	})

	t.Run("returns the expressions correctly", func(t *testing.T) {
		expressionRepositoryMock.
			On("GetAllExpressions", ctx).
			Return([]repositories.Expression{}, nil).
			Once()

		exps, err := expressionService.ListExpressions(ctx)
		require.NoError(t, err)

		assert.Empty(t, exps)

		expectedExps := []repositories.Expression{
			{
				ID:    1,
				Value: "x AND z",
			},
			{
				ID:    2,
				Value: "x OR z",
			},
		}

		expressionRepositoryMock.
			On("GetAllExpressions", ctx).
			Return(expectedExps, nil).
			Once()

		exps, err = expressionService.ListExpressions(ctx)
		require.NoError(t, err)

		assert.Equal(t, expectedExps, exps)
	})
}

func TestExpressionService_UpdateExpression(t *testing.T) {
	expressionRepositoryMock := &repositories.ExpressionRepositoryMock{}
	expressionService := NewExpressionService(WithExpressionRepositoryOption(expressionRepositoryMock))

	ctx := context.Background()

	t.Run("returns error when the expressions is invalid", func(t *testing.T) {
		exp, err := expressionService.UpdateExpression(ctx, &repositories.Expression{
			ID:    0,
			Value: "",
		})

		assert.EqualError(t, err, "invalid expression ID provided")
		assert.Nil(t, exp)

		exp, err = expressionService.CreateExpression(ctx, &repositories.Expression{
			ID:    1,
			Value: "",
		})

		assert.EqualError(t, err, "value can't be empty")
		assert.Nil(t, exp)

		exp, err = expressionService.UpdateExpression(ctx, &repositories.Expression{
			ID:    1,
			Value: "AND",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)

		exp, err = expressionService.UpdateExpression(ctx, &repositories.Expression{
			ID:    1,
			Value: "x AND",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)

		exp, err = expressionService.UpdateExpression(ctx, &repositories.Expression{
			ID:    1,
			Value: "x AND b OR",
		})

		assert.EqualError(t, err, ErrInvalidExpression.Error())
		assert.Nil(t, exp)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		exp := &repositories.Expression{
			ID:    1,
			Value: "x AND y",
		}

		expressionRepositoryMock.
			On("UpdateExpression", ctx, exp).
			Return(nil, errors.New("unexpected error")).
			Once()

		exp, err := expressionService.UpdateExpression(ctx, exp)

		assert.EqualError(t, err, "error updating expression ID 1: unexpected error")
		assert.Nil(t, exp)

		exp = &repositories.Expression{
			ID:    1,
			Value: "x AND y",
		}

		expressionRepositoryMock.
			On("UpdateExpression", ctx, exp).
			Return(nil, repositories.ErrNoRowsAffected).
			Once()

		exp, err = expressionService.UpdateExpression(ctx, exp)

		assert.EqualError(t, err, repositories.ErrNoRowsAffected.Error())
		assert.Nil(t, exp)
	})

	t.Run("updates an expression correctly", func(t *testing.T) {
		exp := &repositories.Expression{
			ID:    1,
			Value: "x AND y",
		}

		expectedExp := &repositories.Expression{
			ID:    1,
			Value: "x AND y",
		}

		expressionRepositoryMock.
			On("UpdateExpression", ctx, exp).
			Return(expectedExp, nil).
			Once()

		exp, err := expressionService.UpdateExpression(ctx, exp)
		require.NoError(t, err)

		assert.Equal(t, expectedExp, exp)
	})

	expressionRepositoryMock.AssertExpectations(t)
}

func TestExpressionService_EvaluateExpression(t *testing.T) {
	expressionRepositoryMock := &repositories.ExpressionRepositoryMock{}
	expressionService := NewExpressionService(WithExpressionRepositoryOption(expressionRepositoryMock))

	ctx := context.Background()

	t.Run("returns error when repository fails", func(t *testing.T) {
		ID := int64(1)

		expressionRepositoryMock.
			On("GetExpressionByID", ctx, ID).
			Return(nil, repositories.ErrExpressionNotFound).
			Once()

		res, err := expressionService.EvaluateExpression(ctx, ID, map[string]int{})

		assert.EqualError(t, err, repositories.ErrExpressionNotFound.Error())
		assert.False(t, res)

		expressionRepositoryMock.
			On("GetExpressionByID", ctx, ID).
			Return(nil, errors.New("unexpected error")).
			Once()

		res, err = expressionService.EvaluateExpression(ctx, ID, map[string]int{})

		assert.EqualError(t, err, "error getting expression ID 1: unexpected error")
		assert.False(t, res)
	})

	t.Run("returns error when a parameter is missing", func(t *testing.T) {
		ID := int64(1)

		expressionRepositoryMock.
			On("GetExpressionByID", ctx, ID).
			Return(&repositories.Expression{
				ID:    ID,
				Value: "x AND z",
			}, nil).
			Once()

		res, err := expressionService.EvaluateExpression(ctx, ID, map[string]int{
			"x": 1,
		})

		assert.EqualError(t, err, `missing parameter "z" for the logical expression "x AND z"`)
		assert.False(t, res)

		expressionRepositoryMock.
			On("GetExpressionByID", ctx, ID).
			Return(&repositories.Expression{
				ID:    ID,
				Value: "(a OR b) AND (c AND d)",
			}, nil).
			Once()

		res, err = expressionService.EvaluateExpression(ctx, ID, map[string]int{
			"b": 1,
			"a": 0,
			"d": 1,
		})

		assert.EqualError(t, err, `missing parameter "c" for the logical expression "(a OR b) AND (c AND d)"`)
		assert.False(t, res)
	})

	t.Run("evaluates the expressions with its parameters", func(t *testing.T) {
		ID := int64(1)

		testCases := []struct {
			expression string
			parameters map[string]int
			expect     bool
		}{
			{
				expression: "a AND z",
				parameters: map[string]int{
					"a": 1,
					"z": 0,
				},
				expect: false,
			},
			{
				expression: "a OR z",
				parameters: map[string]int{
					"a": 1,
					"z": 0,
				},
				expect: true,
			},
			{
				expression: "((x OR y) AND (z OR k) OR j)",
				parameters: map[string]int{
					"x": 1,
					"y": 0,
					"z": 1,
					"k": 0,
					"j": 1,
				},
				expect: true,
			},
			{
				expression: "(x OR y) AND z",
				parameters: map[string]int{
					"x": 1,
					"y": 0,
					"z": 1,
				},
				expect: true,
			},
		}

		for _, tc := range testCases {
			expressionRepositoryMock.
				On("GetExpressionByID", ctx, ID).
				Return(&repositories.Expression{
					ID:    ID,
					Value: tc.expression,
				}, nil).
				Once()

			res, err := expressionService.EvaluateExpression(ctx, ID, tc.parameters)
			require.NoError(t, err)

			assert.Equal(t, tc.expect, res)
		}
	})

	expressionRepositoryMock.AssertExpectations(t)
}
