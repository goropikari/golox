package tlps

// TLPSFunction is struct of tlps function
type TLPSFunction struct {
	declaration   *Function
	closure       *Environment
	IsInitializer bool
}

// NewTLPSFunction is constructor of TLPSFunction
func NewTLPSFunction(declaration *Function, closure *Environment, IsInitializer bool) *TLPSFunction {
	return &TLPSFunction{
		declaration:   declaration,
		closure:       closure,
		IsInitializer: IsInitializer,
	}
}

// Call calls the function
func (lf *TLPSFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(lf.closure)
	for i, param := range lf.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(lf.declaration.Body, environment)
	if err != nil {
		var v interface{} = err
		switch v.(type) {
		case *ReturnValue:
			if lf.IsInitializer {
				return lf.closure.GetAt(0, "this")
			}
			return v.(*ReturnValue).Value, nil
		default:
			return nil, err
		}
	}

	if lf.IsInitializer {
		return lf.closure.GetAt(0, "this")
	}

	return nil, nil
}

// Arity returns arity of function
func (lf *TLPSFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lc *TLPSFunction) Bind(instance *TLPSInstance) *TLPSFunction {
	environment := NewEnvironment(lc.closure)
	environment.Define("this", instance)
	return NewTLPSFunction(lc.declaration, environment, lc.IsInitializer)
}

func (lf *TLPSFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
