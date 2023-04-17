package handlers

type ExpressionRequest struct {
	Expression string `json:"expression" binding:"required"`
}

type CreateExpressionRequest struct {
	ExpressionRequest
}

type UpdateExpressionRequest struct {
	ExpressionRequest
}

type ExpressionResponse struct {
	ID         int64  `json:"id"`
	Expression string `json:"expression"`
}

type CreateExpressionResponse struct {
	ExpressionResponse
}

type ListExpressionsResponse struct {
	ExpressionResponse
}

type UpdateExpressionResponse struct {
	ExpressionResponse
}
