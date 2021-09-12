package tlps

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: methods,
	}
}

func (lc *LoxClass) FindMethod(name string) (*LoxFunction, error) {
	if v, ok := lc.Methods[name]; ok {
		return v, nil
	}

	return nil, nil
}

func (lc *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewLoxInstance(lc)
	initializer, err := lc.FindMethod("init")
	if err != nil {
		return nil, err
	}
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance, nil
}

func (lc *LoxClass) Arity() int {
	initializer, _ := lc.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (lc *LoxClass) String() string {
	return lc.Name
}
