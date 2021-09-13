package golox

// GoLoxFunction is struct of golox function
type GoLoxFunction struct {
	declaration   *Function
	closure       *Environment
	IsInitializer bool
}

// NewGoLoxFunction is constructor of GoLoxFunction
func NewGoLoxFunction(declaration *Function, closure *Environment, IsInitializer bool) *GoLoxFunction {
	return &GoLoxFunction{
		declaration:   declaration,
		closure:       closure,
		IsInitializer: IsInitializer,
	}
}

// Call calls the function
func (lf *GoLoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
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
func (lf *GoLoxFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lc *GoLoxFunction) Bind(instance *GoLoxInstance) *GoLoxFunction {
	environment := NewEnvironment(lc.closure)
	environment.Define("this", instance)
	return NewGoLoxFunction(lc.declaration, environment, lc.IsInitializer)
}

func (lf *GoLoxFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
