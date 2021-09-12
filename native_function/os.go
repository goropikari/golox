package native_function

import (
	"errors"
	"os"
)

// exit(status code)
// ex. exit(1)

// ExitFunc
type ExitFunc struct{}

func NewExitFunc() *ExitFunc {
	return &ExitFunc{}
}

func (ef *ExitFunc) Arity() int {
	return 1
}

func (ef *ExitFunc) Call(arguments []interface{}) (interface{}, error) {
	if len(arguments) == 1 {
		if v, ok := arguments[0].(float64); ok {
			os.Exit(int(v))
		}
		if v, ok := arguments[0].(int); ok {
			os.Exit(v)
		}
	}

	return nil, errors.New("invalid type")
}
