package golox

type GoLoxInstance struct {
	Klass  *GoLoxClass
	Fields map[string]interface{}
}

func NewGoLoxInstance(klass *GoLoxClass) *GoLoxInstance {
	return &GoLoxInstance{
		Klass:  klass,
		Fields: make(map[string]interface{}),
	}
}

func (lc *GoLoxInstance) Get(name *Token) (interface{}, error) {
	if v, ok := lc.Fields[name.Lexeme]; ok {
		return v, nil
	}

	method, err := lc.Klass.FindMethod(name.Lexeme)
	if err != nil {
		return nil, err
	}
	if method != nil {
		return method.Bind(lc), nil
	}

	return nil, RuntimeError.New(name, "Undefied property '"+name.Lexeme+"'.")
}

func (lc *GoLoxInstance) Set(name *Token, value interface{}) {
	lc.Fields[name.Lexeme] = value
}

func (lc *GoLoxInstance) String() string {
	return lc.Klass.Name + " instance"
}
