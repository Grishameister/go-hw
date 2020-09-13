package stack

import "testing"

func TestSliceStack(t *testing.T) {
	stack := SliceStack{}
	if stack.IsEmpty() != true {
		t.Fatal("Stack isnt empty")
	}
	for i := 0; i < 10; i++ {
		stack.Push(rune(i))
		if stack.Size() != i+1 {
			t.Fatal("Push incorrect")
		}
		if stack.Top() != rune(i) {
			t.Fatal("Top incorrect")
		}
	}

	for i := 10; i > 0; i-- {
		if (stack.Size() != i) && stack.Top() != rune(i-1) {
			t.Fatal("Pop incorrect")
		}
		if !stack.Pop() {
			t.Fatal("Pop in not empty stack")
		}
	}
	if stack.Pop() {
		t.Fatal("Pop in empty stack")
	}
}
