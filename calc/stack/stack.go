package stack

type SliceStack []rune

func (s *SliceStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *SliceStack) Size() int {
	return len(*s)
}

func (s *SliceStack) Push(b rune) {
	*s = append(*s, b)
}

func (s *SliceStack) Pop() bool {
	if s.IsEmpty() {
		return false
	}
	*s = (*s)[:len(*s) - 1]
	return true
}

func (s *SliceStack) Top() rune {
	return (*s)[s.Size() - 1]
}

type Node struct {
	value int
	next *Node
}
