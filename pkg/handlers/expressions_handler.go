package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/CaioTeixeira95/logic-exp/pkg/repositories"
	"github.com/CaioTeixeira95/logic-exp/pkg/services"
	"github.com/gin-gonic/gin"
)

type ExpressionHandler struct {
	expressionService services.ExpressionService
}

type ExpressionHandlerOption func(eh *ExpressionHandler)

func (eh *ExpressionHandler) CreateExpression(c *gin.Context) {
	var reqBody CreateExpressionRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		reqErrs := ParseRequestError(err)
		if len(reqErrs) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "request invalid",
				"details": reqErrs["details"],
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Request invalid in some way",
		})
		return
	}

	ctx := c.Request.Context()

	exp, err := eh.expressionService.CreateExpression(ctx, &repositories.Expression{
		Value: reqBody.Expression,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidExpression) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid expression provided",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An Internal Server error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, CreateExpressionResponse{
		ExpressionResponse{
			ID:         exp.ID,
			Expression: exp.Value,
		},
	})
}

func (eh *ExpressionHandler) ListExpressions(c *gin.Context) {
	ctx := c.Request.Context()

	exps, err := eh.expressionService.ListExpressions(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An Internal Server error occurred",
		})
		return
	}

	respBody := make([]ListExpressionsResponse, 0, len(exps))
	for _, exp := range exps {
		respBody = append(respBody, ListExpressionsResponse{
			ExpressionResponse{
				ID:         exp.ID,
				Expression: exp.Value,
			},
		})
	}

	c.JSON(http.StatusOK, respBody)
}

func (eh *ExpressionHandler) UpdateExpression(c *gin.Context) {
	expID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "request invalid",
			"details": "invalid expression ID provided",
		})
		return
	}

	var reqBody UpdateExpressionRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		reqErrs := ParseRequestError(err)
		if len(reqErrs) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "request invalid",
				"details": reqErrs["details"],
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	ctx := c.Request.Context()
	exp, err := eh.expressionService.UpdateExpression(ctx, &repositories.Expression{
		ID:    int64(expID),
		Value: reqBody.Expression,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrNoRowsAffected) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Expression not found",
			})
			return
		}

		if errors.Is(err, services.ErrInvalidExpression) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid expression provided",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An Internal Server error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, UpdateExpressionResponse{
		ExpressionResponse: ExpressionResponse{
			ID:         exp.ID,
			Expression: exp.Value,
		},
	})
}

func (eh *ExpressionHandler) EvaluateExpression(c *gin.Context) {
	expID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "request invalid",
			"details": "invalid expression ID provided",
		})
		return
	}

	parameters := c.Request.URL.Query()

	paramsToEvaluate := make(map[string]int)
	for key, values := range parameters {
		if len(values) != 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("expect exact one value for the key %q but %d were provided", key, len(values)),
			})
			return
		}

		val, err := strconv.Atoi(values[0])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("error converting to integer value %q of the key %q", values[0], key),
			})
			return
		}

		paramsToEvaluate[key] = val
	}

	ctx := c.Request.Context()

	res, err := eh.expressionService.EvaluateExpression(ctx, int64(expID), paramsToEvaluate)
	if err != nil {
		if errors.Is(err, repositories.ErrExpressionNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Expression not found",
			})
			return
		}

		if strings.Contains(err.Error(), "missing parameter") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An Internal Server error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

func NewExpressionHandler(options ...ExpressionHandlerOption) *ExpressionHandler {
	eh := &ExpressionHandler{}

	for _, option := range options {
		option(eh)
	}

	return eh
}

func WithExpressionServiceOption(es services.ExpressionService) ExpressionHandlerOption {
	return func(eh *ExpressionHandler) {
		eh.expressionService = es
	}
}
