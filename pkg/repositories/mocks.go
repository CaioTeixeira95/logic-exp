package repositories

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ExpressionRepositoryMock struct {
	mock.Mock
}

func (er *ExpressionRepositoryMock) CreateExpression(ctx context.Context, exp *Expression) (*Expression, error) {
	args := er.Called(ctx, exp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Expression), args.Error(1)
}

func (er *ExpressionRepositoryMock) GetAllExpressions(ctx context.Context) ([]Expression, error) {
	args := er.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Expression), args.Error(1)
}

func (er *ExpressionRepositoryMock) GetExpressionByID(ctx context.Context, ID int64) (*Expression, error) {
	args := er.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Expression), args.Error(1)
}

func (er *ExpressionRepositoryMock) UpdateExpression(ctx context.Context, exp *Expression) (*Expression, error) {
	args := er.Called(ctx, exp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Expression), args.Error(1)
}

var _ ExpressionRepository = (*ExpressionRepositoryMock)(nil)
