package uniq

import (
	"errors"
	"strconv"
	"strings"
)

type Options struct {
	flags      rune
	field      int
	offset     int
	withoutReg bool
}

type IOpt struct {
	Input  string
	Output string
}

func GenerateCallBack(eq int, value int) func(int) bool {
	if eq == 1 {
		return func(v int) bool {
			return v == value
		}
	} else if eq == 0 {
		return func(v int) bool {
			return v != value
		}
	} else {
		return func(v int) bool {
			return true
		}
	}
}

func FillOptions(args *[]string) (Options, IOpt, error) {
	if args == nil {
		return Options{}, IOpt{}, errors.New("Invalid pointer")
	}
	err := errors.New("usage")
	options := Options{}
	io := IOpt{}

	for i := 1; i < len(*args); i++ {
		arg := (*args)[i]
		switch arg {
		case "-c":
			if options.flags != 0 {
				return Options{}, IOpt{}, err
			}
			options.flags = 'c'
		case "-u":
			if options.flags != 0 {
				return Options{}, IOpt{}, err
			}
			options.flags = 'u'
		case "-d":
			if options.flags != 0 {
				return Options{}, IOpt{}, err
			}
			options.flags = 'd'
		case "-f":
			i++
			if i >= len(*args) {
				return Options{}, IOpt{}, err
			}
			field, err := strconv.Atoi((*args)[i])
			if err != nil {
				return Options{}, IOpt{}, err
			}
			options.field = field
		case "-s":
			i++
			if i >= len(*args) {
				return Options{}, IOpt{}, err
			}
			offset, err := strconv.Atoi((*args)[i])
			if err != nil {
				return Options{}, IOpt{}, err
			}
			options.offset = offset
		case "-i":
			options.withoutReg = true
		default:
			if io.Input == "" {
				io.Input = (*args)[i]
			} else {
				io.Output = (*args)[i]
			}
		}
	}
	return options, io, nil
}

func PrepareForCase(str string) string {
	return strings.ToUpper(str)
}

func PrepareForOffset(str string, num int) string {
	if num <= 0 {
		return str
	}
	i := 0
	position := 0
	for pos, _ := range str {
		if i == num {
			position = pos
			break
		}
		i++
	}
	return str[position:]
}

func PrepareForField(str string, num int) string {
	if num <= 0 {
		return str
	}
	i := 0
	position := 0
	for pos, char := range str {
		if i == num {
			position = pos
			break
		}
		if char == ' ' {
			i++
		}
	}
	return str[position:]
}

type Value struct {
	s     string
	count int
}

func getMapImpl(strings []string, options *Options) map[string]*Value {
	result := make(map[string]*Value, len(strings))

	for i, _ := range strings {
		s := strings[i]
		if options.field > 0 {
			s = PrepareForField(s, options.field)
		}
		if options.offset > 0 {
			s = PrepareForOffset(s, options.offset)
		}
		if options.withoutReg {
			s = PrepareForCase(s)
		}
		_, ok := result[s]
		if ok {
			result[s].count++
		} else {
			result[s] = new(Value)
			result[s].s = strings[i]
			result[s].count = 1
		}
	}
	return result
}

func GetUniqueOrNot(strings []string, options *Options) []string {
	mapOfVals := getMapImpl(strings, options)

	returnValues := make([]string, 0)
	var op rune
	var cb func(int) bool

	switch options.flags {
	case 'u':
		cb = GenerateCallBack(1, 1)
	case 'd':
		cb = GenerateCallBack(0, 1)
	default:
		cb = GenerateCallBack(-1, 1)
	}
	op = options.flags

	for _, v := range mapOfVals {
		if cb(v.count) {
			if op == 'c' {
				v.s = strconv.Itoa(v.count) + " " + v.s
			}
			returnValues = append(returnValues, v.s)
		}
	}
	return returnValues
}
