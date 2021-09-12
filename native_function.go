package tlps

// NativeCallable is interface to call native function
type NativeCallable interface {
	Call([]interface{}) (interface{}, error)
	Arity() int
}

// NativeFunction is struct for native function
type NativeFunction struct {
	Function NativeCallable
}

// NewNativeFunction is constructor of NativeFunction
func NewNativeFunction(function NativeCallable) TLPSCallable {
	return &NativeFunction{
		Function: function,
	}
}

// Call calls native function
func (nf *NativeFunction) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	return nf.Function.Call(args)
}

// Arity returns arity of native function
func (nf *NativeFunction) Arity() int {
	return nf.Function.Arity()
}
