package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CaioTeixeira95/logic-exp/pkg/repositories"
	"github.com/CaioTeixeira95/logic-exp/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpressionHandler_CreateExpression(t *testing.T) {
	er := &repositories.ExpressionRepositoryMock{}
	es := services.NewExpressionService(services.WithExpressionRepositoryOption(er))
	eh := NewExpressionHandler(WithExpressionServiceOption(es))

	t.Setenv("GIN_MODE", gin.TestMode)

	ctx := context.Background()
	endpoint := "/expressions"

	t.Run("returns BadRequest when body is invalid", func(t *testing.T) {
		r := gin.Default()
		r.POST(endpoint, eh.CreateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader("invalid"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "Request invalid in some way"}`, string(respBody))

		req, _ = http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader("{}"))
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp = w.Result()

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "request invalid", "details": {"expression": "this field is required"}}`, string(respBody))

		req, _ = http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(`{"expression": 1}`))
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp = w.Result()

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "request invalid", "details": {"expression": "invalid type provided for this field"}}`, string(respBody))
	})

	t.Run("returns BadRequest when expression is invalid", func(t *testing.T) {
		r := gin.Default()
		r.POST(endpoint, eh.CreateExpression)

		reqBody := `
			{
				"expression": "(x AND z"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "Invalid expression provided"}`, string(respBody))
	})

	t.Run("returns InternalServerError when an unexpected error occurs", func(t *testing.T) {
		r := gin.Default()
		r.POST(endpoint, eh.CreateExpression)

		reqBody := `
			{
				"expression": "(x AND z)"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		er.
			On("CreateExpression", req.Context(), &repositories.Expression{
				Value: "(x AND z)",
			}).
			Return(nil, errors.New("unexpected error")).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, `{"error":"An Internal Server error occurred"}`, string(respBody))
	})

	t.Run("creates a new expression successfully", func(t *testing.T) {
		r := gin.Default()
		r.POST(endpoint, eh.CreateExpression)

		reqBody := `
			{
				"expression": "(x AND z)"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		er.
			On("CreateExpression", req.Context(), &repositories.Expression{
				Value: "(x AND z)",
			}).
			Return(&repositories.Expression{
				ID:    1,
				Value: "(x AND z)",
			}, nil).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.JSONEq(t, `{"id":1, "expression":"(x AND z)"}`, string(respBody))
	})

	er.AssertExpectations(t)
}

func TestExpressionHandler_ListExpressions(t *testing.T) {
	er := &repositories.ExpressionRepositoryMock{}
	es := services.NewExpressionService(services.WithExpressionRepositoryOption(er))
	eh := NewExpressionHandler(WithExpressionServiceOption(es))

	t.Setenv("GIN_MODE", gin.TestMode)

	ctx := context.Background()
	endpoint := "/expressions"

	t.Run("returns InternalServer error when an unexpected error occurs", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.ListExpressions)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()

		er.
			On("GetAllExpressions", req.Context()).
			Return(nil, errors.New("unexpected error")).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, `{"error": "An Internal Server error occurred"}`, string(respBody))
	})

	t.Run("lists all expressions successfully", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.ListExpressions)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()

		er.
			On("GetAllExpressions", req.Context()).
			Return([]repositories.Expression{}, nil).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, `[]`, string(respBody))

		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		w = httptest.NewRecorder()

		er.
			On("GetAllExpressions", req.Context()).
			Return([]repositories.Expression{
				{
					ID:    1,
					Value: "x AND z",
				},
				{
					ID:    2,
					Value: "(x AND z) OR y",
				},
				{
					ID:    3,
					Value: "(x OR b OR (a AND z))",
				},
			}, nil).
			Once()

		r.ServeHTTP(w, req)

		resp = w.Result()

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		wantsBody := `
			[
				{
					"id": 1,
					"expression": "x AND z"
				},
				{
					"id": 2,
					"expression": "(x AND z) OR y"
				},
				{
					"id": 3,
					"expression": "(x OR b OR (a AND z))"
				}
			]
		`

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, wantsBody, string(respBody))
	})

	er.AssertExpectations(t)
}

func TestExpressionHandler_UpdateExpression(t *testing.T) {
	er := &repositories.ExpressionRepositoryMock{}
	es := services.NewExpressionService(services.WithExpressionRepositoryOption(er))
	eh := NewExpressionHandler(WithExpressionServiceOption(es))

	t.Setenv("GIN_MODE", gin.TestMode)

	ctx := context.Background()
	endpoint := "/expressions/:id"

	t.Run("returns BadRequest when body is invalid", func(t *testing.T) {
		r := gin.Default()
		r.PUT(endpoint, eh.UpdateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader("invalid"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"message":"invalid request"}`, string(respBody))

		req, _ = http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader("{}"))
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp = w.Result()

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "request invalid", "details": {"expression": "this field is required"}}`, string(respBody))

		req, _ = http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader(`{"expression": 1}`))
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp = w.Result()

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "request invalid", "details": {"expression": "invalid type provided for this field"}}`, string(respBody))
	})

	t.Run("returns BadRequest when expression is invalid", func(t *testing.T) {
		r := gin.Default()
		r.PUT(endpoint, eh.UpdateExpression)

		reqBody := `
			{
				"expression": "(x AND z"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "Invalid expression provided"}`, string(respBody))
	})

	t.Run("returns NotFound when expression is not found", func(t *testing.T) {
		r := gin.Default()
		r.PUT(endpoint, eh.UpdateExpression)

		reqBody := `
			{
				"expression": "(x AND z)"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		er.
			On("UpdateExpression", req.Context(), &repositories.Expression{
				ID:    1,
				Value: "(x AND z)",
			}).
			Return(nil, repositories.ErrNoRowsAffected).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.JSONEq(t, `{"error": "Expression not found"}`, string(respBody))
	})

	t.Run("returns InternalServerError when an unexpected error occurs", func(t *testing.T) {
		r := gin.Default()
		r.PUT(endpoint, eh.UpdateExpression)

		reqBody := `
			{
				"expression": "(x AND z)"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		er.
			On("UpdateExpression", req.Context(), &repositories.Expression{
				ID:    1,
				Value: "(x AND z)",
			}).
			Return(nil, errors.New("unexpected error")).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, `{"error":"An Internal Server error occurred"}`, string(respBody))
	})

	t.Run("updates an expression successfully", func(t *testing.T) {
		r := gin.Default()
		r.PUT(endpoint, eh.UpdateExpression)

		reqBody := `
			{
				"expression": "(x AND z)"
			}
		`
		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/expressions/1", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		er.
			On("UpdateExpression", req.Context(), &repositories.Expression{
				ID:    1,
				Value: "(x AND z)",
			}).
			Return(&repositories.Expression{
				ID:    1,
				Value: "(x AND z)",
			}, nil).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, `{"id":1, "expression":"(x AND z)"}`, string(respBody))
	})

	er.AssertExpectations(t)
}

func TestExpressionHandler_EvaluateExpression(t *testing.T) {
	er := &repositories.ExpressionRepositoryMock{}
	es := services.NewExpressionService(services.WithExpressionRepositoryOption(er))
	eh := NewExpressionHandler(WithExpressionServiceOption(es))

	t.Setenv("GIN_MODE", gin.TestMode)

	ctx := context.Background()
	endpoint := "/evaluate/:id"

	t.Run("returns BadRequest error when has multiple values for the same query param", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/evaluate/1?x=1&x=2", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "expect exact one value for the key \"x\" but 2 were provided"}`, string(respBody))
	})

	t.Run("returns BadRequest error when query param isn't integer", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/evaluate/1?x=abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error":"error converting to integer value \"abc\" of the key \"x\""}`, string(respBody))
	})

	t.Run("returns NotFound when expression is not found", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/evaluate/1?x=1", nil)
		w := httptest.NewRecorder()

		er.
			On("GetExpressionByID", req.Context(), int64(1)).
			Return(nil, repositories.ErrExpressionNotFound).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.JSONEq(t, `{"error": "Expression not found"}`, string(respBody))
	})

	t.Run("returns BadRequest when a parameter is missing", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/evaluate/1?x=1", nil)
		w := httptest.NewRecorder()

		er.
			On("GetExpressionByID", req.Context(), int64(1)).
			Return(&repositories.Expression{
				ID:    1,
				Value: "(x AND y)",
			}, nil).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.JSONEq(t, `{"error": "missing parameter \"y\" for the logical expression \"(x AND y)\""}`, string(respBody))
	})

	t.Run("returns InternalServerError when an unexpected error occurs", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/evaluate/1?x=1&y=0", nil)
		w := httptest.NewRecorder()

		er.
			On("GetExpressionByID", req.Context(), int64(1)).
			Return(nil, errors.New("unexpected error")).
			Once()

		r.ServeHTTP(w, req)

		resp := w.Result()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, `{"error": "An Internal Server error occurred"}`, string(respBody))
	})

	t.Run("evaluates expressions successfully", func(t *testing.T) {
		r := gin.Default()
		r.GET(endpoint, eh.EvaluateExpression)

		testCases := []struct {
			expID                             int64
			expression, parameters, wantsBody string
		}{
			{
				expID:      1,
				expression: "x AND z",
				parameters: "x=1&z=0",
				wantsBody:  `{"result": false}`,
			},
			{
				expID:      2,
				expression: "((x OR y) AND (z OR k) OR j)",
				parameters: "x=1&y=0&z=1&k=0&j=1",
				wantsBody:  `{"result": true}`,
			},
			{
				expID:      3,
				expression: "(x OR y) AND z",
				parameters: "x=1&y=0&z=1",
				wantsBody:  `{"result": true}`,
			},
			{
				expID:      4,
				expression: "x AND x",
				parameters: "x=1",
				wantsBody:  `{"result": true}`,
			},
		}

		for _, tc := range testCases {
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/evaluate/%d?%s", tc.expID, tc.parameters), nil)
			w := httptest.NewRecorder()

			er.
				On("GetExpressionByID", req.Context(), tc.expID).
				Return(&repositories.Expression{
					ID:    tc.expID,
					Value: tc.expression,
				}, nil).
				Once()

			r.ServeHTTP(w, req)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.JSONEq(t, tc.wantsBody, string(respBody))
		}
	})

	er.AssertExpectations(t)
}
