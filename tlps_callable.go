package tlps

// TLPSCallable is interface
type TLPSCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, error)
	Arity() int
}
