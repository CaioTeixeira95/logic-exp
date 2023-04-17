package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogicalExpressionParameters(t *testing.T) {
	testCases := []struct {
		expression string
		expect     LogicalExpressionParametersSet
	}{
		{
			expression: "(x OR y) AND z",
			expect: LogicalExpressionParametersSet{
				"x": struct{}{},
				"y": struct{}{},
				"z": struct{}{},
			},
		},
		{
			expression: "((x OR y) AND z)",
			expect: LogicalExpressionParametersSet{
				"x": struct{}{},
				"y": struct{}{},
				"z": struct{}{},
			},
		},
		{
			expression: "x AND y OR z",
			expect: LogicalExpressionParametersSet{
				"x": struct{}{},
				"y": struct{}{},
				"z": struct{}{},
			},
		},
		{
			expression: "x AND y OR z OR (x OR y)",
			expect: LogicalExpressionParametersSet{
				"x": struct{}{},
				"y": struct{}{},
				"z": struct{}{},
			},
		},
		{
			expression: "a OR b AND (c OR (b AND f))",
			expect: LogicalExpressionParametersSet{
				"a": struct{}{},
				"b": struct{}{},
				"c": struct{}{},
				"f": struct{}{},
			},
		},
		{
			expression: "x AND",
			expect: LogicalExpressionParametersSet{
				"x": struct{}{},
			},
		},
		{
			expression: "xyz AND z OR y AND x",
			expect: LogicalExpressionParametersSet{
				"xyz": struct{}{},
				"x":   struct{}{},
				"y":   struct{}{},
				"z":   struct{}{},
			},
		},
		{
			expression: "AND OR",
			expect:     LogicalExpressionParametersSet{},
		},
	}

	for _, tc := range testCases {
		got := GetLogicalExpressionParameters(tc.expression)
		assert.Equal(t, tc.expect, got)
	}
}

func TestIsLogicalExpressionValid(t *testing.T) {
	testCases := []struct {
		expression string
		expect     bool
	}{
		{
			expression: "x + 1",
			expect:     false,
		},
		{
			expression: "x AND z",
			expect:     true,
		},
		{
			expression: "(y OR z) AND x OR a AND f",
			expect:     true,
		},
		{
			expression: "OR",
			expect:     false,
		},
		{
			expression: "(x OR )",
			expect:     false,
		},
	}

	for _, tc := range testCases {
		isValid := IsLogicalExpressionValid(tc.expression)
		assert.Equal(t, tc.expect, isValid)
	}
}

func TestEvaluateExpression(t *testing.T) {
	testCases := []struct {
		expression string
		parameters map[string]int
		expect     bool
		err        string
	}{
		{
			expression: "a AND z",
			parameters: map[string]int{
				"a": 1,
				"z": 0,
			},
			expect: false,
			err:    "",
		},
		{
			expression: "((x OR y) AND z)",
			parameters: map[string]int{
				"x": 1,
				"y": 0,
				"z": 1,
			},
			expect: true,
			err:    "",
		},
		{
			expression: "((x OR y) AND z)",
			parameters: map[string]int{
				"x": 0,
				"y": 0,
				"z": 1,
			},
			expect: false,
			err:    "",
		},
		{
			expression: "((x OR y) AND z)",
			parameters: map[string]int{
				"x": 0,
				"y": 1,
				"z": 0,
			},
			expect: false,
			err:    "",
		},
		{
			expression: "(x AND z) OR (a OR b) AND y",
			parameters: map[string]int{
				"x": 0,
				"y": 1,
				"z": 0,
				"a": 0,
				"b": 1,
			},
			expect: true,
			err:    "",
		},
		{
			expression: "(x AND x)",
			parameters: map[string]int{
				"x": 1,
			},
			expect: true,
			err:    "",
		},
		// Expressions with error
		{
			expression: "x AND",
			parameters: map[string]int{
				"x": 1,
			},
			expect: false,
			err:    "error creating evaluable expression: Unexpected end of expression",
		},
		{
			expression: "(x AND x",
			parameters: map[string]int{
				"x": 1,
			},
			expect: false,
			err:    "error creating evaluable expression: Unbalanced parenthesis",
		},
		{
			expression: "AND",
			parameters: map[string]int{
				"x": 1,
			},
			expect: false,
			err:    "error creating evaluable expression: Cannot transition token types from UNKNOWN [<nil>] to LOGICALOP [&&]",
		},
		// Missing or Extra param
		{
			expression: "z AND y",
			parameters: map[string]int{
				"z": 1,
			},
			expect: false,
			err:    `error evaluating expression "z && y" with parameters map[z:1]: No parameter 'y' found.`,
		},
		{
			expression: "z AND y",
			parameters: map[string]int{
				"z": 1,
				"y": 1,
				"x": 0,
			},
			expect: true,
			err:    "",
		},
		// Expressions that doesn't result in boolean results
		{
			expression: "x + z",
			parameters: map[string]int{
				"x": 2,
				"z": 1,
			},
			expect: false,
			err:    `error evaluating expression "x + z" with parameters map[x:2 z:1]: Value 'true' cannot be used with the modifier '+', it is not a number`,
		},
		{
			expression: "x OR z OR 1",
			parameters: map[string]int{
				"x": 0,
				"z": 0,
			},
			expect: false,
			err:    "error evaluating expression \"x || z || 1\" with parameters map[x:0 z:0]: Value '1' cannot be used with the logical operator '||', it is not a bool",
		},
	}

	for _, tc := range testCases {
		res, err := EvaluateLogicalExpression(tc.expression, tc.parameters)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err)
			assert.False(t, res)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, res)
		}
	}
}
