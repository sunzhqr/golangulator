package parser

import (
	"github.com/Knetic/govaluate"
)

func Eval(expr string) (float64, error) {
	e, err := govaluate.NewEvaluableExpression(expr)
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
