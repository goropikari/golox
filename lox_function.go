package tlps

// LoxFunction is struct of lox function
type LoxFunction struct {
	declaration   *Function
	closure       *Environment
	IsInitializer bool
}

// NewLoxFunction is constructor of LoxFunction
func NewLoxFunction(declaration *Function, closure *Environment, IsInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		IsInitializer: IsInitializer,
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
func (lf *LoxFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lc *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(lc.closure)
	environment.Define("this", instance)
	return NewLoxFunction(lc.declaration, environment, lc.IsInitializer)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
