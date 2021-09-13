package golox

// GoLoxCallable is interface
type GoLoxCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, error)
	Arity() int
}
