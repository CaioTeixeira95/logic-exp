package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

// pattern is a regex that aims get only the operands of a logical expression.
const pattern = `(?m)[a-z]+`

var regex *regexp.Regexp

type LogicalExpressionParametersSet map[string]struct{}

func init() {
	regex = regexp.MustCompile(pattern)
}

// GetLogicalExpressionParameters returns the operands of a logical expression as a set.
func GetLogicalExpressionParameters(logicalExpression string) LogicalExpressionParametersSet {
	parameters := regex.FindAllString(logicalExpression, -1)

	parametersSet := make(LogicalExpressionParametersSet)
	for _, param := range parameters {
		parametersSet[param] = struct{}{}
	}

	return parametersSet
}

// IsLogicalExpressionValid validates whether a logical expression is valid or not.
// Only boolean expressions are accept.
func IsLogicalExpressionValid(logicalExpression string) bool {
	// Validate if the expressions is a boolean expression and not a
	// mathematic expression, for instance.
	parametersSet := GetLogicalExpressionParameters(logicalExpression)

	parameters := make(map[string]int, len(parametersSet))
	for key := range parametersSet {
		parameters[key] = 0
	}

	_, err := EvaluateLogicalExpression(logicalExpression, parameters)

	return err == nil
}

// EvaluateLogicalExpression replace the operands in a logical expression with the parameters
// passed by parameter on it and evaluate the result of the expression.
func EvaluateLogicalExpression(logicalExpression string, parameters map[string]int) (bool, error) {
	logicalExpression = strings.ReplaceAll(logicalExpression, "AND", "&&")
	logicalExpression = strings.ReplaceAll(logicalExpression, "OR", "||")

	exp, err := govaluate.NewEvaluableExpression(logicalExpression)
	if err != nil {
		return false, fmt.Errorf("error creating evaluable expression: %w", err)
	}

	paramsToEvaluate := make(map[string]interface{}, len(parameters))
	for key, value := range parameters {
		paramsToEvaluate[key] = value > 0
	}

	result, err := exp.Evaluate(paramsToEvaluate)
	if err != nil {
		return false, fmt.Errorf("error evaluating expression %q with parameters %v: %w", exp.String(), parameters, err)
	}

	return result.(bool), nil
}
