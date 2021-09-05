package tlps

type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
}

func NewEnvironment(environment *Environment) *Environment {
	return &Environment{
		Values:    make(map[string]interface{}, 0),
		Enclosing: environment,
	}
}

func (e *Environment) Get(name *Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, RuntimeError.New(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(name *Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return RuntimeError.New(name, "Undefined variable '"+name.Lexeme+"'.")
}
