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
	length, err := rpn.RPN(slice)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rpn.Calculate(length))
}
