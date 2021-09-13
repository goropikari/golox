package golox

type GoLoxClass struct {
	Name       string
	Superclass *GoLoxClass
	Methods    map[string]*GoLoxFunction
}

func NewGoLoxClass(name string, superclass *GoLoxClass, methods map[string]*GoLoxFunction) *GoLoxClass {
	return &GoLoxClass{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (lc *GoLoxClass) FindMethod(name string) (*GoLoxFunction, error) {
	if v, ok := lc.Methods[name]; ok {
		return v, nil
	}

	if lc.Superclass != nil {
		return lc.Superclass.FindMethod(name)
	}

	return nil, nil
}

func (lc *GoLoxClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewGoLoxInstance(lc)
	initializer, err := lc.FindMethod("init")
	if err != nil {
		return nil, err
	}
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance, nil
}

func (lc *GoLoxClass) Arity() int {
	initializer, _ := lc.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (lc *GoLoxClass) String() string {
	return lc.Name
}
