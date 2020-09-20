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

func isValidBranches(expr []rune) bool {
	branches := stack.Stack{}

	for _, v := range expr {
		if v == '(' {
			branches.Push(v)
		}
		if v == ')' {
			if branches.Len() == 0 {
				return false
			}
			branches.Pop()
		}
	}
	if branches.Len() != 0 {
		return false
	}
	return true
}

func isValidSymbols(expr []rune) bool {
	if len(expr) == 0 {
		return false
	}
	if isOptionRune(expr[0]) || isOptionRune(expr[len(expr)-1]) {
		return false
	}

	if !isOptionRune(expr[0]) && !isBranch(expr[0]) && !unicode.IsDigit(expr[0]) {
		return false
	}

	prev := expr[0]
	for i := 1; i < len(expr); i++ {
		if !isOptionRune(expr[i]) && !isBranch(expr[i]) && !unicode.IsDigit(expr[i]) {
			return false
		}
		if isOptionRune(prev) && isOptionRune(expr[i]) {
			return false
		}
		if !isBranch(expr[i]) {
			prev = expr[i]
		}
	}
	return true
}

func RPN(expr []rune) ([]string, error) {
	invalidExpr := errors.New("Invalid")
	if !isValidBranches(expr) {
		return make([]string, 0), invalidExpr
	}
	if !isValidSymbols(expr) {
		return make([]string, 0), invalidExpr
	}
	var isNumber bool
	operations := stack.Stack{}
	rpn := make([]string, 0)
	number := make([]rune, 0)
	for _, v := range expr {
		if unicode.IsDigit(v) {
			if !isNumber {
				isNumber = true
			}
			number = append(number, v)
		} else {
			if isNumber {
				rpn = append(rpn, string(number))
				number = number[:0]
			}
			isNumber = false
		}
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
	if isNumber {
		rpn = append(rpn, string(number))
	}
	for operations.Len() != 0 {
		rpn = append(rpn, string(operations.Peek().(int32)))
		operations.Pop()
	}
	return rpn, nil
}

func Calculate(expr []string) int32 {
	s := stack.Stack{}
	for _, v := range expr {
		if isOptionString(v) {
			op1 := s.Peek().(int32)
			s.Pop()
			op2 := s.Peek().(int32)
			s.Pop()
			var result int32
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
			value, _ := strconv.Atoi(v)
			s.Push(rune(value))
		}
	}
	return s.Peek().(int32)
}
