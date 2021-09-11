package tlps

// LoxCallable is interface
type LoxCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, error)
	Arity() int
}
