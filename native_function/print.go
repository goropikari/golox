package native_function

import "fmt"

type PrintFunc struct{}

func NewPrintFunc() *PrintFunc {
	return &PrintFunc{}
}

func (pf *PrintFunc) Arity() int {
	return -1
}

func (pf *PrintFunc) Call(arguments []interface{}) (interface{}, error) {
	fmt.Print(arguments...)
	return nil, nil
}
