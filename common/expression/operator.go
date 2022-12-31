package expression

import (
	"fmt"
	"github.com/pkg/errors"
)

// Precedence returns > 0 if op1 > op2, or < 0 if op1 < op2, otherwise 0
type Precedence func(op1, op2 string) int

type Operator struct {
	op    string
	left  interface{}
	right interface{}
}

func (o *Operator) String() string {
	var l, r string

	switch tl := o.left.(type) {
	case variable:
		l = tl.String()
	case int64:
		l = fmt.Sprintf("%d", tl)
	case *Operator:
		l = tl.String()
	}

	switch tr := o.right.(type) {
	case variable:
		r = tr.String()
	case int64:
		r = fmt.Sprintf("%d", tr)
	case *Operator:
		r = tr.String()
	}

	return fmt.Sprintf("(%s %s %s)", l, o.op, r)
}

func (o *Operator) EvalKnown(vars map[string]int64) (int64, error) {
	var l, r int64
	var el, er error
	switch tl := o.left.(type) {
	case variable:
		l, el = tl.EvalKnown(vars)
		if el == nil {
			o.left = l
		}
	case int64:
		l = tl
	case *Operator:
		l, el = tl.EvalKnown(vars)
		if el == nil {
			o.left = l
		}
	}
	switch tr := o.right.(type) {
	case variable:
		r, er = tr.EvalKnown(vars)
		if er == nil {
			o.right = r
		}
	case int64:
		r = tr
	case *Operator:
		r, er = tr.EvalKnown(vars)
		if er == nil {
			o.right = r
		}
	}

	if el != nil || er != nil {
		return 0, errors.New("unable to eval operator")
	}

	switch o.op {
	case "-":
		return l - r, nil
	case "+":
		return l + r, nil
	case "*":
		return l * r, nil
	case "/":
		return l / r, nil
	}
	panic(errors.New("unknown operator"))
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
