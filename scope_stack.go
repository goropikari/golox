package tlps

import "github.com/goropikari/tlps/collections/stack"

// ScopeStack is struct of stack for scopes
type ScopeStack struct {
	Stack *stack.Stack
}

// NewScopeStack is constructor of ScopeStack
func NewScopeStack() *ScopeStack {
	return &ScopeStack{stack.NewStack()}
}

// Push adds an item in stack
func (s *ScopeStack) Push(x map[string]bool) {
	s.Stack.Push(x)
}

// Pop pops an item from stack
func (s *ScopeStack) Pop() map[string]bool {
	return s.Stack.Pop().(map[string]bool)
}

// Peek returns top item in stack, and don't modity the stack.
func (s *ScopeStack) Peek() map[string]bool {
	return s.Stack.Peek().(map[string]bool)
}

// IsEmpty checks that stack is empty
func (s *ScopeStack) IsEmpty() bool {
	return s.Stack.IsEmpty()
}

// Size returns stack size
func (s *ScopeStack) Size() int {
	return s.Stack.Size()
}

// Get returns i th element from top
func (s *ScopeStack) Get(i int) (map[string]bool, error) {
	m, err := s.Stack.Get(i)
	if err != nil {
		return nil, err
	}
	return m.(map[string]bool), err
}
