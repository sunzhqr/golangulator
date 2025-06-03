package usecase

import (
	"github.com/sunzhqr/golangulator/internal/domain"
	"github.com/sunzhqr/golangulator/pkg/parser"
)

type calculatorUseCase struct{}

func NewCalculatorUseCase() domain.Calculator {
	return &calculatorUseCase{}
}

func (c *calculatorUseCase) Eval(expr string) (float64, error) {
	return parser.Eval(expr)
}
