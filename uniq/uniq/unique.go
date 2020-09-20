package uniq

import (
	"strconv"
	"strings"
)

type Options struct {
	Flags      rune
	Field      int
	Offset     int
	WithoutReg bool
}

type IOpt struct {
	Input  string
	Output string
}

func GenerateCallBack(eq int8, value int) func(int) bool {
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
	result := make(map[string]*Value, 0)

	for i, _ := range strings {
		s := strings[i]
		if options.Field > 0 {
			s = PrepareForField(s, options.Field)
		}
		if options.Offset > 0 {
			s = PrepareForOffset(s, options.Offset)
		}
		if options.WithoutReg {
			s = PrepareForCase(s)
		}
		if _, ok := result[s]; ok {
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

	switch options.Flags {
	case 'u':
		cb = GenerateCallBack(1, 1)
	case 'd':
		cb = GenerateCallBack(0, 1)
	default:
		cb = GenerateCallBack(-1, 1)
	}
	op = options.Flags

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
