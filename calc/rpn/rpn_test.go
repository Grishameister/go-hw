package rpn

import "testing"

func TestIsOptionSuccess(t *testing.T) {
	input := []string{"+", "-", "/", "*"}
	for _, v := range input {
		out := isOption(v)
		if !out {
			t.Fatalf("Check failed on : %s", input)
		}
	}

}

func TestIsOptionFail(t *testing.T) {
	input := "23245"
	out := isOption(input)

	if out {
		t.Fatalf("Check failed on : %s", input)
	}
}

func TestIsValidExprSuccess(t *testing.T) {
	input := []rune("(((9+2)*2-3))")

	out := isValidExpr(input)

	if !out {
		t.Fatal("Check Success failed on isValid")
	}

}

func TestIsValidExprFailLeft(t *testing.T) {
	input := []rune("((9+2)*2-3))")

	out := isValidExpr(input)

	if out {
		t.Fatal("Check Fail failed on isValid")
	}
}

func TestIsValidExprFailRight(t *testing.T) {
	input := []rune("((9+2)*2-3")

	out := isValidExpr(input)

	if out {
		t.Fatal("Check Fail failed on isValid")
	}
}

func TestRPNSuccess(t *testing.T) {
	input := []rune("(12+25*3-4/5)+9")
	out, err := RPN(input)
	if err != nil {
		t.Fatal("Invalid expression")
	}
	expected := []string{"12", "25", "3", "*", "+", "4", "5", "/", "-", "9", "+"}

	if len(out) != len(expected) {
		t.Fatal("Lengths mismatch")
	}
	for i := 0; i < len(expected); i++ {
		if out[i] != expected[i] {
			t.Fatal("strings aren't equal")
		}
	}
}

func TestRPNBFail(t *testing.T) {
	input := []rune("(12+25*3-4/5)+9)")

	_, err := RPN(input)
	if err == nil {
		t.Fatal("Invalid expression")
	}
}

func TestCalculateSuccess(t *testing.T) {
	input := []string{"12", "25", "3", "*", "+", "4", "5", "/", "-", "9", "+"}

	out := Calculate(input)
	if out != 95 {
		t.Fatal("Bad Calculate")
	}
}
