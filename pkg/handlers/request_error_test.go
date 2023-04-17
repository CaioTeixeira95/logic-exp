package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

func TestParseRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	var (
		ts  TestStruct
		err error
	)

	ctx.Request, err = http.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"test":"test"}`))
	require.NoError(t, err)

	err = ctx.ShouldBindJSON(&ts)
	require.Error(t, err)

	reqError := ParseRequestError(err)

	assert.Equal(t, gin.H{
		"details": gin.H{
			"name": "this field is required",
			"age":  "this field is required",
		},
	}, reqError)

	ctx.Request, err = http.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"name": 1123}`))
	require.NoError(t, err)

	err = ctx.ShouldBindJSON(&ts)
	require.Error(t, err)

	reqError = ParseRequestError(err)

	assert.Equal(t, gin.H{
		"details": gin.H{
			"name": "invalid type provided for this field",
		},
	}, reqError)
}
