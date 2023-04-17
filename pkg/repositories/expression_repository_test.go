package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeExpressionFixture(t *testing.T, ctx context.Context, expression string) (ID int64) {
	query := `
		INSERT INTO expressions (expression) VALUES ($1) RETURNING id
	`
	err := testConn.QueryRowContext(ctx, query, expression).Scan(&ID)
	require.NoError(t, err)

	return
}

func deleteExpressionsFixture(t *testing.T, ctx context.Context) {
	query := `
		DELETE FROM expressions
	`
	_, err := testConn.ExecContext(ctx, query)
	require.NoError(t, err)
}

func TestDefaultRepository_CreateExpression(t *testing.T) {
	er := NewRepository(WithDatabaseOption(testConn))

	ctx := context.Background()

	t.Run("insert expressions successfully", func(t *testing.T) {
		exp, err := er.CreateExpression(ctx, &Expression{
			Value: "x AND z",
		})
		require.NoError(t, err)

		assert.NotNil(t, exp)
		assert.NotEmpty(t, exp.ID)
	})
}

func TestDefaultRepository_UpdateExpression(t *testing.T) {
	er := NewRepository(WithDatabaseOption(testConn))

	ctx := context.Background()

	t.Run("returns error when no rows is affected", func(t *testing.T) {
		_, err := er.UpdateExpression(ctx, &Expression{
			Value: "x AND z",
		})

		assert.EqualError(t, err, ErrNoRowsAffected.Error())

		_, err = er.UpdateExpression(ctx, &Expression{
			ID:    999,
			Value: "x AND z",
		})

		assert.EqualError(t, err, ErrNoRowsAffected.Error())
	})

	t.Run("updates expressions successfully", func(t *testing.T) {
		ID := makeExpressionFixture(t, ctx, "x AND z")

		exp, err := er.UpdateExpression(ctx, &Expression{
			ID:    ID,
			Value: "(x AND z OR y)",
		})
		require.NoError(t, err)

		assert.NotEqual(t, exp.Value, "x AND z")
	})
}

func TestDefaultRepository_GetAllExpressions(t *testing.T) {
	er := NewRepository(WithDatabaseOption(testConn))

	ctx := context.Background()

	t.Run("returns all expressions successfully", func(t *testing.T) {
		deleteExpressionsFixture(t, ctx)

		exps, err := er.GetAllExpressions(ctx)
		require.NoError(t, err)

		assert.Empty(t, exps)

		expID1 := makeExpressionFixture(t, ctx, "x AND y")
		expID2 := makeExpressionFixture(t, ctx, "x OR y")
		expID3 := makeExpressionFixture(t, ctx, "a AND b")

		expectedExps := []Expression{
			{
				ID:    expID1,
				Value: "x AND y",
			},
			{
				ID:    expID2,
				Value: "x OR y",
			},
			{
				ID:    expID3,
				Value: "a AND b",
			},
		}

		exps, err = er.GetAllExpressions(ctx)
		require.NoError(t, err)

		assert.Equal(t, expectedExps, exps)
	})
}

func TestDefaultRepository_GetExpressionByID(t *testing.T) {
	er := NewRepository(WithDatabaseOption(testConn))

	ctx := context.Background()

	t.Run("returns error when expressions is not found", func(t *testing.T) {
		deleteExpressionsFixture(t, ctx)

		exp, err := er.GetExpressionByID(ctx, 1)

		assert.EqualError(t, err, ErrExpressionNotFound.Error())
		assert.Nil(t, exp)
	})

	t.Run("returns error when expressions is not found", func(t *testing.T) {
		expID := makeExpressionFixture(t, ctx, "x AND y")

		expectedExp := &Expression{
			ID:    expID,
			Value: "x AND y",
		}

		exp, err := er.GetExpressionByID(ctx, expID)
		require.NoError(t, err)

		assert.Equal(t, expectedExp, exp)
	})
}
