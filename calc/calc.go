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
	for scanner.Scan() {
		expression = scanner.Text()
		break
	}
	slice := []rune(expression)
	length, _ := rpn.RPN(slice)
	fmt.Println(rpn.Calculate(length))
}
