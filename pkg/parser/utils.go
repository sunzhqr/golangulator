package parser

import (
	"fmt"
	"regexp"
	"strings"
)

func preprocessExpression(expr string) string {
	expr = strings.ReplaceAll(expr, " ", "")

	rePower := regexp.MustCompile(`(\d+(?:\.\d+)?|\([^()]*\))\^(\d+(?:\.\d+)?|\([^()]*\))`)
	for {
		loc := rePower.FindStringSubmatchIndex(expr)
		if loc == nil {
			break
		}
		base := expr[loc[2]:loc[3]]
		exp := expr[loc[4]:loc[5]]
		powExpr := fmt.Sprintf("pow(%s,%s)", base, exp)
		expr = expr[:loc[0]] + powExpr + expr[loc[1]:]
	}

	rePercent := regexp.MustCompile(`(\d+(?:\.\d+)?)%`)
	expr = rePercent.ReplaceAllString(expr, "($1/100)")

	return expr
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func translateParseError(err error) string {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "Unclosed function"):
		return "не закрыта скобка у функции"
	case strings.Contains(msg, "Unrecognized character"):
		return "использован недопустимый символ"
	case strings.Contains(msg, "Expected closing bracket"):
		return "отсутствует закрывающая скобка"
	case strings.Contains(msg, "Missing operand"):
		return "в выражении отсутствует число или оператор"
	default:
		return "проверьте корректность выражения"
	}
}

func translateEvalError(err error) string {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "divide by zero"):
		return "деление на ноль невозможно"
	case strings.Contains(msg, "invalid parameter"):
		return "ошибка в параметрах функции"
	default:
		return "не удалось выполнить вычисление"
	}
}
