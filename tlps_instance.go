package tlps

type TLPSInstance struct {
	Klass  *TLPSClass
	Fields map[string]interface{}
}

func NewTLPSInstance(klass *TLPSClass) *TLPSInstance {
	return &TLPSInstance{
		Klass:  klass,
		Fields: make(map[string]interface{}),
	}
}

func (lc *TLPSInstance) Get(name *Token) (interface{}, error) {
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

func (lc *TLPSInstance) Set(name *Token, value interface{}) {
	lc.Fields[name.Lexeme] = value
}

func (lc *TLPSInstance) String() string {
	return lc.Klass.Name + " instance"
}
