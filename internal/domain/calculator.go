package domain

type Calculator interface {
	Eval(expr string) (float64, error)
}
