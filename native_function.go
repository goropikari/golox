package tlps

import "time"

type ClockFunc struct{}

func NewClockFunc() *ClockFunc {
	return &ClockFunc{}
}

func (cf *ClockFunc) Arity() int {
	return 0
}

func (cf *ClockFunc) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (cf *ClockFunc) String() string {
	return "<native fn>"
}
