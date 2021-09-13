package golox

import "github.com/goropikari/golox/collections/stack"

// IndentStack is struct of stack for storing indentation level
type IndentStack struct {
	Stack *stack.Stack
}

// NewIndentStack is constructor of IndentStack
func NewIndentStack() *IndentStack {
	return &IndentStack{stack.NewStack()}
}

// Push adds an item in stack
func (s *IndentStack) Push(x int) {
	s.Stack.Push(x)
}

// Pop pops an item from stack
func (s *IndentStack) Pop() int {
	return s.Stack.Pop().(int)
}

// Peek returns top item in stack, and don't modity the stack.
func (s *IndentStack) Peek() int {
	return s.Stack.Peek().(int)
}

// IsEmpty checks that stack is empty
func (s *IndentStack) IsEmpty() bool {
	return s.Stack.IsEmpty()
}

// Size returns stack size
func (s *IndentStack) Size() int {
	return s.Stack.Size()
}
