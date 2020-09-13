package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"uniq/uniq"
)

func open(io string, def *os.File, flag int) (*os.File, error) {
	descr := def
	if io != "" {
		file, err := os.OpenFile(io, flag, 0644)
		if err != nil {
			return descr, errors.New(io)
		}
		descr = file
	}
	return descr, nil
}

func main() {
	options, io, err := uniq.FillOptions(&os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	in, err := open(io.Input, os.Stdin, os.O_RDONLY)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if in != os.Stdin {
			err := in.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	out, err := open(io.Output, os.Stdout, os.O_WRONLY)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if out != os.Stdout {
			err := out.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	scanner := bufio.NewScanner(in)
	printer := bufio.NewWriter(out)

	strings := make([]string, 0)
	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	for _, v := range uniq.GetUniqueOrNot(strings, &options) {
		printer.WriteString(v)
		printer.WriteRune('\n')
	}
	printer.Flush()
}
