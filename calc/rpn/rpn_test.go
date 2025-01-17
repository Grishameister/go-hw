package rpn

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsOptionSuccess(t *testing.T) {
	input := []string{"+", "-", "/", "*"}
	for _, v := range input {
		assert.Equal(t, isOptionString(v), true)
	}
}

func TestIsOptionFail(t *testing.T) {
	input := "23245"
	assert.Equal(t, isOptionString(input), false)
}

func TestIsValidSymbolsSuccess(t *testing.T) {
	input := []rune("(((9.2+2)*2-3))")
	assert.Equal(t, isValidExpr(input), true)
}

func TestIsValidSymbolsFail(t *testing.T) {
	tests := []struct {
		input []rune
		out   bool
	}{
		{
			input: []rune("1++1"),
			out:   false,
		},
		{
			input: []rune("9+2*2-3-"),
			out:   false,
		},
		{
			input: []rune("9+2-a*2-3-"),
			out:   false,
		},
		{
			input: []rune("a+2-a*2-3"),
			out:   false,
		},
		{
			input: []rune("2-a*2-3-a"),
			out:   false,
		},
		{
			input: []rune{},
			out:   false,
		},
		{
			input: []rune("((9+2)*2-3+4))"),
			out:   false,
		},
		{
			input: []rune("((9+2)*2-3"),
			out:   false,
		},
		{
			input: []rune(")(9+2)*2-3"),
			out:   false,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			out := isValidExpr(tt.input)
			assert.Equal(t, out, tt.out)
		})
	}
}

func TestRPNSuccess(t *testing.T) {
	input := []rune("(12+25*3-4/5)+9")
	out, err := RPN(input)
	if err != nil {
		t.Fatal("Invalid expression")
	}
	expected := []string{"12", "25", "3", "*", "+", "4", "5", "/", "-", "9", "+"}
	assert.Equal(t, expected, out)
}

func TestRPNFail(t *testing.T) {
	tests := []struct {
		input []rune
		out   []string
	}{
		{
			input: []rune("(12+25*3-4/5)+9)"),
			out:   []string{},
		},
		{
			input: []rune("(12+25*3-4/5)-+9"),
			out:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			out, err := RPN(tt.input)
			require.Error(t, err)
			assert.Equal(t, out, tt.out)
		})
	}
}

func TestCalculateSuccess(t *testing.T) {
	input := []string{"12", "25", "3", "*", "+", "4", "5", "/", "-", "9", "+"}
	out, err := Calculate(input)
	assert.Equal(t, out, 94.75)
	require.NoError(t, err)
}

func TestCalculateFail(t *testing.T) {
	input := []string{"12", "25.2.2", "3", "*", "+", "4", "5", "/", "-", "9", "+"}
	out, err := Calculate(input)
	assert.Equal(t, out, float64(0))
	require.Error(t, err)
}
