package tlps

type TLPSClass struct {
	Name       string
	Superclass *TLPSClass
	Methods    map[string]*TLPSFunction
}

func NewTLPSClass(name string, superclass *TLPSClass, methods map[string]*TLPSFunction) *TLPSClass {
	return &TLPSClass{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (lc *TLPSClass) FindMethod(name string) (*TLPSFunction, error) {
	if v, ok := lc.Methods[name]; ok {
		return v, nil
	}

	if lc.Superclass != nil {
		return lc.Superclass.FindMethod(name)
	}

	return nil, nil
}

func (lc *TLPSClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewTLPSInstance(lc)
	initializer, err := lc.FindMethod("init")
	if err != nil {
		return nil, err
	}
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance, nil
}

func (lc *TLPSClass) Arity() int {
	initializer, _ := lc.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (lc *TLPSClass) String() string {
	return lc.Name
}
