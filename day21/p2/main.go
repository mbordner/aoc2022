package main

import (
	"fmt"
	"github.com/mbordner/aoc2022/common/expression"
	"github.com/mbordner/aoc2022/common/file"
	"regexp"
	"strconv"
)

var (
	reExpr  = regexp.MustCompile(`^(\w+):\s(\w+)\s(\+|\-|\*|\/)\s(\w+)$`)
	reValue = regexp.MustCompile(`^(\w+):\s(\-?\d+)$`)
)

type ExpressionMap map[string]Expression
type ValueMap map[string]int64

func (vm ValueMap) Has(s string) bool {
	if _, e := vm[s]; e {
		return true
	}
	return false
}

func (em ExpressionMap) Has(s string) bool {
	if _, e := em[s]; e {
		return true
	}
	return false
}

func main() {
	exprMap, valMap := getExpressions("../data.txt")
	if len(exprMap) > 0 && len(valMap) > 0 {
		if exprMap.Has("root") {

			delete(valMap, "humn")
			r := exprMap["root"]
			exprMap["root"] = Expression{
				name:  r.name,
				left:  r.left,
				op:    "-",
				right: r.right,
			}

			expr, err := expression.NewParser(exprMap["root"].String(exprMap, valMap))
			if err != nil {
				panic(err)
			}

			_, _ = expr.EvalKnown(valMap)

			fmt.Println("0 = ", expr.String())

			leftExpr := exprMap[exprMap["root"].left].String(exprMap, valMap)
			rightExpr := exprMap[exprMap["root"].right].String(exprMap, valMap)

			lExpr, err := expression.NewParser(leftExpr)
			if err != nil {
				panic(err)
			}

			rExpr, err := expression.NewParser(rightExpr)
			if err != nil {
				panic(err)
			}

			_, lErr := lExpr.EvalKnown(valMap)
			_, rErr := rExpr.EvalKnown(valMap)

			var operator *expression.Operator

			if lErr != nil {
				_, operator, err = lExpr.RootOperator().InverseOperationToVariableExpression(rExpr.RootOperator())
			} else if rErr != nil {
				_, operator, err = rExpr.RootOperator().InverseOperationToVariableExpression(lExpr.RootOperator())
			}

			if err != nil {
				panic(err)
			}

			fmt.Println("humn = ", operator.String())

			val, _ := operator.EvalKnown(valMap)

			fmt.Println("humn = ", val)

		} else {
			panic("no root expression")
		}
	}
}

type Expression struct {
	name  string
	left  string
	op    string
	right string
}

func (e Expression) String(em ExpressionMap, vm ValueMap) string {
	var left, right string

	if em.Has(e.left) {
		left = em[e.left].String(em, vm)
	} else if vm.Has(e.left) {
		left = fmt.Sprintf("%d", vm[e.left])
	} else {
		left = e.left
	}

	if em.Has(e.right) {
		right = em[e.right].String(em, vm)
	} else if vm.Has(e.right) {
		right = fmt.Sprintf("%d", vm[e.right])
	} else {
		right = e.right
	}

	return fmt.Sprintf(`(%s %s %s)`, left, e.op, right)
}

func getExpressions(path string) (ExpressionMap, ValueMap) {
	lines, _ := file.GetLines(path)
	exprMap := make(ExpressionMap)
	valMap := make(ValueMap)

	for _, line := range lines {
		if reValue.MatchString(line) {
			matches := reValue.FindStringSubmatch(line)
			var value int64
			value, _ = strconv.ParseInt(matches[2], 10, 64)
			valMap[matches[1]] = value
		} else if reExpr.MatchString(line) {
			matches := reExpr.FindStringSubmatch(line)
			expression := Expression{
				name:  matches[1],
				left:  matches[2],
				op:    matches[3],
				right: matches[4],
			}
			exprMap[matches[1]] = expression
		}
	}

	return exprMap, valMap
}
