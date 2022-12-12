package expression

import "github.com/pkg/errors"

// Precedence returns > 0 if op1 > op2, or < 0 if op1 < op2, otherwise 0
type Precedence func(op1, op2 string) int

type Operator struct {
	op    string
	left  interface{}
	right interface{}
}

func (o *Operator) Eval(vars map[string]int64) int64 {
	var l, r int64
	switch tl := o.left.(type) {
	case variable:
		l = tl.Eval(vars)
	case int64:
		l = tl
	case *Operator:
		l = tl.Eval(vars)
	}
	switch tr := o.right.(type) {
	case variable:
		r = tr.Eval(vars)
	case int64:
		r = tr
	case *Operator:
		r = tr.Eval(vars)
	}
	switch o.op {
	case "-":
		return l - r
	case "+":
		return l + r
	case "*":
		return l * r
	case "/":
		return l / r
	}
	panic(errors.New("unknown operator"))
}

func CompareOperator(op1, op2 string) int {
	return 0
}

func IsBinary(op string) bool {
	return true
}
