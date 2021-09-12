package tlps

type LoxInstance struct {
	Klass  *LoxClass
	Fields map[string]interface{}
}

func NewLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		Klass:  klass,
		Fields: make(map[string]interface{}),
	}
}

func (lc *LoxInstance) Get(name *Token) (interface{}, error) {
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

func (lc *LoxInstance) Set(name *Token, value interface{}) {
	lc.Fields[name.Lexeme] = value
}

func (lc *LoxInstance) String() string {
	return lc.Klass.Name + " instance"
}
