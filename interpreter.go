package tlps

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	runtime *Runtime
}

func NewInterpreter(runtime *Runtime) *Interpreter {
	return &Interpreter{runtime: runtime}
}

func (i *Interpreter) Interpret(expression Expr) (string, error) {
	value, err := i.evaluate(expression)
	if err != nil {
		i.runtime.RuntimeError(err)
	}
	// fmt.Println(stringfy(value))

	return stringfy(value), err
}

func (i *Interpreter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Greater:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case Less:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BangEqual:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) != right.(float64), nil
	case EqualEqual:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) == right.(float64), nil
	case Minus:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case Plus:
		if isFloat64(left) && isFloat64(right) {
			return left.(float64) + right.(float64), nil
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		}

		return nil, RuntimeError.New(expr.Operator, "Operands must be two numbers or two strings.")
	case Slash:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case Star:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	// Unreachable.
	return nil, RuntimeError.New(nil, "Unreachable")
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Bang:
		return !i.isTruthy(right), nil
	case Minus:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}

	// Unreachable
	return nil, RuntimeError.New(nil, "Unreachable")
}

func checkNumberOperand(operator *Token, operand interface{}) error {
	if reflect.ValueOf(operand).Kind() == reflect.Float64 {
		return nil
	}
	return RuntimeError.New(operator, "Operand must be a number.")
}

func checkNumberOperands(operator *Token, left interface{}, right interface{}) error {
	if reflect.ValueOf(left).Kind() == reflect.Float64 &&
		reflect.ValueOf(right).Kind() == reflect.Float64 {
		return nil
	}
	return RuntimeError.New(operator, "Operands must be a number.")
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch v := reflect.ValueOf(object); v.Kind() {
	case reflect.Bool:
		return object.(bool)
	}

	return true
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func isFloat64(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Float64
}

func isString(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.String
}

func stringfy(object interface{}) string {
	if object == nil {
		return "nil"
	}

	return fmt.Sprint(object)
}
