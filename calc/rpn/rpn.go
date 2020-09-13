package rpn

import (
	"calc/stack"
	"errors"
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

func isOption(s string) bool {
	return s == "+" || s == "-" || s == "/" || s == "*"
}

func isValidExpr(expr []rune) bool {
	branches := stack.SliceStack{}

	for _, v := range expr {
		if v == '(' {
			branches.Push(v)
		}
		if v == ')' {
			if branches.IsEmpty() {
				return false
			}
			branches.Pop()
		}
	}
	if !branches.IsEmpty() {
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
	operations := stack.SliceStack{}
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
			for !operations.IsEmpty() && v != '(' {
				op := operations.Top()
				p := ops[op]
				if p >= priority {
					rpn = append(rpn, string(op))
					operations.Pop()
				} else {
					operations.Push(v)
					break
				}
			}
			if operations.IsEmpty() || v == '(' {
				operations.Push(v)
			}
		}
		if v == ')' {
			op := operations.Top()
			for op != '(' {
				rpn = append(rpn, string(op))
				operations.Pop()
				op = operations.Top()
			}
			operations.Pop()
		}
	}
	if isNumber {
		rpn = append(rpn, string(number))
	}
	for !operations.IsEmpty() {
		rpn = append(rpn, string(operations.Top()))
		operations.Pop()
	}
	return rpn, nil
}

func Calculate(expr []string) int {
	s := stack.SliceStack{}
	for _, v := range expr {
		if isOption(v) {
			op1 := s.Top()
			s.Pop()
			op2 := s.Top()
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
	return int(s.Top())
}
