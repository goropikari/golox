package stack

import "errors"

// Stack is a stack data structure
type Stack struct {
	Data   []interface{}
	Length int
}

// NewStack is constructor of Stack
func NewStack() *Stack {
	return new(Stack)
}

// Push adds an item in stack
func (s *Stack) Push(x interface{}) {
	s.Data = append(s.Data, x)
	s.Length++
}

// Pop pops an item from stack
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return -1
	}

	item := s.Peek()
	s.Length--
	s.Data = s.Data[:s.Length]
	return item
}

// Peek returns top item in stack, and don't modity the stack.
func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return -1
	}
	return s.Data[s.Length-1]
}

// IsEmpty checks that stack is empty
func (s *Stack) IsEmpty() bool {
	return s.Length == 0
}

// Size returns stack size
func (s *Stack) Size() int {
	return s.Length
}

// Get returns i th element of Stack from top
func (s *Stack) Get(i int) (interface{}, error) {
	if i < 0 || i >= s.Size() {
		return nil, errors.New("BoundError")
	}

	return s.Data[s.Size()-1-i], nil
}
