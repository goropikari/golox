package golox

// Environment is struct of environment
type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
}

// NewEnvironment is constructor of Environment
func NewEnvironment(environment *Environment) *Environment {
	return &Environment{
		Values:    make(map[string]interface{}, 0),
		Enclosing: environment,
	}
}

// Get returns value associated with name
func (e *Environment) Get(name *Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, RuntimeError.New(name, "Undefined variable '"+name.Lexeme+"'.")
}

// GetAt return value associated with name at distance depth environment.
func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	return e.Ancestor(distance).Values[name], nil
}

// Ancestor returns `distance` th environment
func (e *Environment) Ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}

	return environment
}

// Define defines variable
func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

// Assign assigns value
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

// AssignAt assigns value at `distance` th environment
func (e *Environment) AssignAt(distance int, name *Token, value interface{}) error {
	e.Ancestor(distance).Values[name.Lexeme] = value

	return nil
}
