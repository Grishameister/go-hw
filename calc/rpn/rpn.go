package rpn

import (
	"errors"
	"github.com/golang-collections/collections/stack"
	"strconv"
	"unicode"
)

var ops = map[rune]int{
	'(': 1,
	'+': 2,
	'-': 2,
	'*': 3,
	'/': 3,
}

func isOptionString(s string) bool {
	return s == "+" || s == "-" || s == "/" || s == "*"
}

func isBranch(sym rune) bool {
	return sym == '(' || sym == ')'
}

func isOptionRune(sym rune) bool {
	return sym == '+' || sym == '-' || sym == '/' || sym == '*'
}

func isPoint(sym rune) bool {
	return sym == '.'
}

func isValidExpr(expr []rune) bool {
	branches := stack.Stack{}
	if len(expr) == 0 {
		return false
	}
	if isOptionRune(expr[0]) || isOptionRune(expr[len(expr)-1]) || isPoint(expr[0]) || isPoint(expr[len(expr)-1]) {
		return false
	}

	if !isOptionRune(expr[0]) && expr[0] == ')' && !unicode.IsDigit(expr[0]) {
		return false
	}

	if isBranch(expr[0]) {
		branches.Push(expr[0])
	}

	prev := expr[0]
	for i := 1; i < len(expr); i++ {
		if !isOptionRune(expr[i]) && !isBranch(expr[i]) && !unicode.IsDigit(expr[i]) && !isPoint(expr[i]) {
			return false
		}
		if (isOptionRune(prev) || isPoint(prev)) && (isOptionRune(expr[i]) || isPoint(expr[i])) {
			return false
		}
		if !isBranch(expr[i]) {
			prev = expr[i]
		} else {
			if expr[i] == '(' {
				branches.Push(expr[i])
			}
			if expr[i] == ')' {
				if branches.Len() == 0 {
					return false
				}
				branches.Pop()
			}
		}
	}
	if branches.Len() != 0 {
		return false
	}
	return true
}

func RPN(expr []rune) ([]string, error) {
	invalidExpr := errors.New("Invalid")
	if !isValidExpr(expr) {
		return make([]string, 0), invalidExpr
	}
	var isNumber bool
	operations := stack.Stack{}
	rpn := make([]string, 0)
	number := make([]rune, 0)
	var empty []rune
	for _, v := range expr {
		if unicode.IsDigit(v) || isPoint(v) {
			if !isNumber {
				isNumber = true
			}
			number = append(number, v)
		} else {
			if isNumber {
				rpn = append(rpn, string(number))
				number = empty
			}
			isNumber = false

			if priority, err := ops[v]; err {
				for operations.Len() != 0 && v != '(' {
					op := operations.Peek().(rune)
					p := ops[op]
					if p >= priority {
						rpn = append(rpn, string(op))
						operations.Pop()
					} else {
						operations.Push(v)
						break
					}
				}
				if operations.Len() == 0 || v == '(' {
					operations.Push(v)
				}
			}
			if v == ')' {
				op := operations.Peek().(rune)
				for op != '(' {
					rpn = append(rpn, string(op))
					operations.Pop()
					op = operations.Peek().(rune)
				}
				operations.Pop()
			}
		}
	}
	if isNumber {
		rpn = append(rpn, string(number))
	}
	for operations.Len() != 0 {
		rpn = append(rpn, string(operations.Peek().(int32)))
		operations.Pop()
	}
	return rpn, nil
}

func Calculate(expr []string) (float64, error) {
	s := stack.Stack{}
	for _, v := range expr {
		if isOptionString(v) {
			op1 := s.Peek().(float64)
			s.Pop()
			op2 := s.Peek().(float64)
			s.Pop()
			var result float64
			switch v {
			case "+":
				result = op1 + op2
			case "-":
				result = op2 - op1
			case "/":
				result = op1 / op2
			case "*":
				result = op1 * op2
			}
			s.Push(result)
		} else {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, err
			}
			s.Push(value)
		}
	}
	return s.Peek().(float64), nil
}
