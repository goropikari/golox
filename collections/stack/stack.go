package stack

// Stack is a stack data structure
type Stack struct {
	data   []int
	length int
}

// NewStack is constructor of Stack
func NewStack() *Stack {
	return new(Stack)
}

// Push adds an item in stack
func (s *Stack) Push(x int) {
	s.data = append(s.data, x)
	s.length++
}

// Pop pops an item from stack
func (s *Stack) Pop() int {
	if s.IsEmpty() {
		return -1
	}

	item := s.Top()
	s.length--
	s.data = s.data[:s.length]
	return item
}

// Top returns top item in stack, and don't modity the stack.
func (s *Stack) Top() int {
	if s.IsEmpty() {
		return -1
	}
	return s.data[s.length-1]
}

// IsEmpty checks that stack is empty
func (s *Stack) IsEmpty() bool {
	return s.length == 0
}

// Size returns stack size
func (s *Stack) Size() int {
	return s.length
}
