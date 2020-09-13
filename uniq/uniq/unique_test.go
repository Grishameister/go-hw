package uniq

import (
	"fmt"
	"testing"
)

func TestFillOptionsSuccess(t *testing.T) {
	input := []string{"unique", "-d", "-i", "-f", "2", "-s", "1", "input.txt", "output.txt"}
	expectedOpt := Options{
		flags:      'd',
		field:      2,
		offset:     1,
		withoutReg: true,
	}
	expectedIO := IOpt{
		Input:  "input.txt",
		Output: "output.txt",
	}

	outOpt, outIO, err := FillOptions(&input)

	if outOpt != expectedOpt {
		fmt.Println(outOpt, expectedOpt)
		t.Fatal("Fill Options doesn't parse")
	}
	if outIO != expectedIO {
		t.Fatal("Fill IO doesn't parse")
	}
	if err != nil {
		t.Fatal("Error happens")
	}
}

func TestFillOptionsDoubleFlag(t *testing.T) {
	input := []string{"unique", "-c", "-d"}
	input2 := []string{"unique", "-d", "-u"}
	input3 := []string{"unique", "-u", "-c"}

	_, _, err := FillOptions(&input)
	if err == nil {
		t.Fatal("Error happens")
	}

	_, _, err2 := FillOptions(&input2)
	if err2 == nil {
		t.Fatal("Error happens")
	}

	_, _, err3 := FillOptions(&input3)
	if err3 == nil {
		t.Fatal("Error happens")
	}
}

func TestFillOptionsConvertF(t *testing.T) {
	input := []string{"unique", "-f", "-d"}
	_, _, err := FillOptions(&input)
	if err == nil {
		t.Fatal("Error happens")
	}
}

func TestFillOptionsConvertS(t *testing.T) {
	input := []string{"unique", "-s", "abcde"}
	_, _, err := FillOptions(&input)
	if err == nil {
		t.Fatal("Error happens")
	}
}

func TestFillOptionsLastFS(t *testing.T) {
	input := []string{"unique", "-s"}
	input2 := []string{"unique", "-f"}

	_, _, err := FillOptions(&input)
	if err == nil {
		t.Fatal("Error happens")
	}

	_, _, err2 := FillOptions(&input2)
	if err2 == nil {
		t.Fatal("Error happens")
	}
}

func TestFillOptionsNil(t *testing.T) {
	var input *[]string
	_, _, err := FillOptions(input)
	if err == nil {
		t.Fatal("Error happens")
	}
}

func TestPrepareForCase(t *testing.T) {
	input := "AbCd"
	out := PrepareForCase(input)
	expected := "ABCD"

	if out != expected {
		t.Fatal("Case error")
	}
}

func TestPrepareForOffset(t *testing.T) {
	input := "AbCd"
	out := PrepareForOffset(input, 2)
	expected := "Cd"

	if out != expected {
		t.Fatal("Offset error")
	}
	out2 := PrepareForOffset(input, -1)
	if out2 != input {
		t.Fatal("Offset error")
	}
}

func TestPrepareForField(t *testing.T) {
	input := "AbCd like me"
	out1 := PrepareForField(input, 2)
	expected1 := "me"

	if out1 != expected1 {
		t.Fatal("Field error")
	}

	out2 := PrepareForField(input, -1)
	if out2 != input {
		t.Fatal("Field error")
	}
}

func TestGetUniqueOrNot(t *testing.T) {
	input := []string{"I LOVE MUSIC.",
		"I love music.",
		"I LoVe MuSiC.",
		"I love MuSIC of Kartik.",
		"I love music of kartik.",
		"Thanks.",
	}

	opt := Options{
		flags:      'c',
		field:      2,
		offset:     1,
		withoutReg: true,
	}

	out := GetUniqueOrNot(input, &opt)
	expected := []string{"3 I LOVE MUSIC.", "2 I love MuSIC of Kartik.", "1 Thanks."}

	if len(out) != len(expected) {
		t.Fatal("Unique doesn't work")
	}
	for i := 0; i < len(expected); i++ {
		in := false
		for j := 0; j < len(out); j++ {
			if expected[i] == expected[j] {
				in = true
			}
		}
		if !in {
			t.Fatal("Mismatch strings")
		}
	}
}

func TestGetUniqueOrNotD(t *testing.T) {
	input := []string{"I LOVE MUSIC.",
		"I love music.",
		"I LoVe MuSiC and not.",
		"I love MuSIC of Kartik.",
		"I love music of kartik.",
		"Thanks.",
	}

	opt := Options{
		flags:      'd',
		field:      1,
		offset:     1,
		withoutReg: true,
	}

	out := GetUniqueOrNot(input, &opt)
	expected := []string{"I LOVE MUSIC.", "I love MuSIC of Kartik."}

	if len(out) != len(expected) {
		t.Fatal("Unique doesn't work")
	}
	for i := 0; i < len(expected); i++ {
		in := false
		for j := 0; j < len(out); j++ {
			if expected[i] == expected[j] {
				in = true
			}
		}
		if !in {
			t.Fatal("Mismatch strings")
		}
	}
}

func TestGetUniqueOrNotU(t *testing.T) {
	input := []string{"I LOVE MUSIC.",
		"I love music.",
		"I LoVe MuSiC and not.",
		"I love MuSIC of Kartik.",
		"I love music of kartik.",
		"Thanks.",
	}

	opt := Options{
		flags:      'u',
		field:      2,
		offset:     1,
		withoutReg: true,
	}

	out := GetUniqueOrNot(input, &opt)
	expected := []string{"I LoVe MuSiC and not.", "Thanks."}

	if len(out) != len(expected) {
		t.Fatal("Unique doesn't work")
	}
	for i := 0; i < len(expected); i++ {
		in := false
		for j := 0; j < len(out); j++ {
			if expected[i] == expected[j] {
				in = true
			}
		}
		if !in {
			t.Fatal("Mismatch strings")
		}
	}
}
