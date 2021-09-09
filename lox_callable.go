package tlps

type LoxCallable interface {
	Call(*Interpreter, []interface{}) (interface{}, error)
	Arity() int
}
