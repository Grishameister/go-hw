package uniq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareForCase(t *testing.T) {
	input := "AbCd"
	assert.Equal(t, PrepareForCase(input), "ABCD")
}

func TestPrepareForOffset(t *testing.T) {
	input := "AbCd"
	assert.Equal(t, PrepareForOffset(input, 2), "Cd")
	assert.Equal(t, PrepareForOffset(input, -1), input)
}

func TestPrepareForField(t *testing.T) {
	input := "AbCd like me"
	assert.Equal(t, PrepareForField(input, 2), "me")
	assert.Equal(t, PrepareForField(input, -1), input)
}

func TestGetUniqueOrNot(t *testing.T) {
	tests := []struct {
		input    []string
		opt      Options
		expected []string
	}{
		{
			input: []string{"I LOVE MUSIC.",
				"I love music.",
				"I LoVe MuSiC.",
				"I love MuSIC of Kartik.",
				"I love music of kartik.",
				"Thanks.",
			},
			opt: Options{
				Flags:      'c',
				Field:      2,
				Offset:     1,
				WithoutReg: true,
			},
			expected: []string{"3 I LOVE MUSIC.",
				"2 I love MuSIC of Kartik.",
				"1 Thanks.",
			},
		},
		{
			input: []string{"I LOVE MUSIC.",
				"I love music.",
				"I LoVe MuSiC and not.",
				"I love MuSIC of Kartik.",
				"I love music of kartik.",
				"Thanks.",
			},
			opt: Options{
				Flags:      'd',
				Field:      1,
				Offset:     1,
				WithoutReg: true,
			},
			expected: []string{
				"I LOVE MUSIC.",
				"I love MuSIC of Kartik.",
			},
		},
		{
			input: []string{"I LOVE MUSIC.",
				"I love music.",
				"I LoVe MuSiC and not.",
				"I love MuSIC of Kartik.",
				"I love music of kartik.",
				"Thanks.",
			},
			opt: Options{
				Flags:      'u',
				Field:      2,
				Offset:     1,
				WithoutReg: true,
			},
			expected: []string{
				"I LoVe MuSiC and not.",
				"Thanks.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.opt.Flags), func(t *testing.T) {
			expectedMap := map[string]bool{}
			outMap := map[string]bool{}

			for _, v := range GetUniqueOrNot(tt.input, &tt.opt) {
				outMap[v] = true
			}
			for _, v := range tt.expected {
				expectedMap[v] = true
			}
			assert.Equal(t, outMap, expectedMap)
		})
	}
}
