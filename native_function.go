package tlps

import "time"

// ClockFunc is struct of clock function
type ClockFunc struct{}

// NewClockFunc is constructor of ClockFunc
func NewClockFunc() *ClockFunc {
	return &ClockFunc{}
}

// Arity returns 0
func (cf *ClockFunc) Arity() int {
	return 0
}

// Call return now time
func (cf *ClockFunc) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (cf *ClockFunc) String() string {
	return "<native fn>"
}
