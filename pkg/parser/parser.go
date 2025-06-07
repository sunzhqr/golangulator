package parser

import (
	"errors"
	"fmt"
	"math"

	"github.com/Knetic/govaluate"
)

func Eval(expr string) (float64, error) {
	expr = preprocessExpression(expr)

	functions := map[string]govaluate.ExpressionFunction{
		"pow": func(args ...interface{}) (interface{}, error) {
			base, ok1 := toFloat(args[0])
			exp, ok2 := toFloat(args[1])
			if !ok1 || !ok2 {
				return nil, errors.New("некорректные аргументы для функции pow")
			}
			return math.Pow(base, exp), nil
		},
	}

	e, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		return 0, fmt.Errorf("ошибка в выражении: %s", translateParseError(err))
	}

	result, err := e.Evaluate(map[string]interface{}{}) // не передаём nil!
	if err != nil {
		return 0, fmt.Errorf("ошибка при вычислении: %s", translateEvalError(err))
	}

	val, ok := toFloat(result)
	if !ok {
		return 0, errors.New("результат имеет неподдерживаемый формат")
	}

	if math.IsInf(val, 0) || math.IsNaN(val) {
		return 0, errors.New("деление на ноль невозможно или результат не определён")
	}

	return val, nil
}
