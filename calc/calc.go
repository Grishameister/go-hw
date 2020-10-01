package main

import (
	"bufio"
	"calc/rpn"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var expression string
	if scanner.Scan() {
		expression = scanner.Text()
	}
	slice := []rune(expression)
	expr, err := rpn.RPN(slice)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if value, err := rpn.Calculate(expr); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}
}
