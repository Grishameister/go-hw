package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"uniq/uniq"
)

func open(inout string, def *os.File, flag int) (*os.File, error) {
	descr := def
	if inout != "" {
		file, err := os.OpenFile(inout, flag, 0644)
		if err != nil {
			return descr, errors.New(inout)
		}
		descr = file
	}
	return descr, nil
}

func FillOpt() (uniq.Options, uniq.IOpt, error) {
	err := errors.New("usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
	cFlag := flag.Bool("c", false, "With Counters")
	uFlag := flag.Bool("u", false, "Unique")
	dFlag := flag.Bool("d", false, "Not Unique")

	iFlag := flag.Bool("i", false, "Register")

	var field int
	var offset int
	flag.IntVar(&field, "f", 0, "Fields")
	flag.IntVar(&offset, "s", 0, "Offset")

	opt := uniq.Options{}

	flag.Parse()

	if *cFlag {
		opt.Flags = 'c'
	}
	if *uFlag {
		if opt.Flags != 0 {
			return uniq.Options{}, uniq.IOpt{}, err
		}
		opt.Flags = 'u'
	}
	if *dFlag {
		if opt.Flags != 0 {
			return uniq.Options{}, uniq.IOpt{}, err
		}
		opt.Flags = 'd'
	}

	if *iFlag {
		opt.WithoutReg = true
	}

	opt.Field = field
	opt.Offset = offset

	iopt := uniq.IOpt{}

	length := len(flag.Args())

	if length > 2 {
		return uniq.Options{}, uniq.IOpt{}, err
	}

	if length > 0 {
		iopt.Input = flag.Arg(0)
	}
	if length > 1 {
		iopt.Output = flag.Arg(1)
	}

	return opt, iopt, nil
}

func uniqueWork() error {
	options, io, err := FillOpt()
	if err != nil {

		return err
	}

	in, err := open(io.Input, os.Stdin, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer func() {
		if in != os.Stdin {
			err := in.Close()
			if err != nil {
				return
			}
		}
	}()

	out, err := open(io.Output, os.Stdout, os.O_WRONLY)
	if err != nil {
		return err
	}
	defer func() {
		if out != os.Stdout {
			err := out.Close()
			if err != nil {
				return
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
		if _, err := printer.WriteString(v + "\n"); err != nil {
			fmt.Println(err)
			return err
		}
	}
	if err := printer.Flush(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	if err := uniqueWork(); err != nil {
		fmt.Println(err)
	}
}
