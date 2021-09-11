package tlps

// LoxFunction is struct of lox function
type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

// NewLoxFunction is constructor of LoxFunction
func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

// Call calls the function
func (lf *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(lf.closure)
	for i, param := range lf.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(lf.declaration.Body, environment)
	if err != nil {
		var v interface{} = err
		switch v.(type) {
		case *ReturnValue:
			return v.(*ReturnValue).Value, nil
		default:
			return nil, err
		}
	}

	return nil, nil
}

// Arity returns arity of function
func (lf *LoxFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
