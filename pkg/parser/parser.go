package parser

import (
	"math"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

func Eval(expr string) (float64, error) {
	expr = preprocessExpression(expr)

	functions := map[string]govaluate.ExpressionFunction{
		"pow": func(args ...interface{}) (interface{}, error) {
			base := args[0].(float64)
			exp := args[1].(float64)
			return math.Pow(base, exp), nil
		},
	}

	e, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		return 0, err
	}
	result, err := e.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	if val, ok := result.(float64); ok {
		return val, nil
	}
	return 0, nil
}

func preprocessExpression(expr string) string {
	expr = strings.ReplaceAll(expr, " ", "")

	re := regexp.MustCompile(`(\d+(?:\.\d+)?|\([\d\+\-\*/\.^]+\))\^(\d+(?:\.\d+)?|\([\d\+\-\*/\.^]+\))`)
	for {
		loc := re.FindStringSubmatchIndex(expr)
		if loc == nil {
			break
		}
		a := expr[loc[2]:loc[3]]
		b := expr[loc[4]:loc[5]]
		expr = expr[:loc[0]] + "pow(" + a + "," + b + ")" + expr[loc[1]:]
	}

	rePercent := regexp.MustCompile(`(\d+(?:\.\d+)?)%`)
	expr = rePercent.ReplaceAllString(expr, "($1/100)")

	return expr
}
