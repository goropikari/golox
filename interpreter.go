package tlps

import (
	"fmt"
	"reflect"
)

// Interpreter is struct of interpreter
type Interpreter struct {
	Runtime *Runtime
}

// NewInterpreter is constructor of Interpreter
func NewInterpreter(runtime *Runtime) *Interpreter {
	globals := runtime.Globals

	globals.Define("clock", NewClockFunc())

	return &Interpreter{
		Runtime: runtime,
	}
}

// Interpret interprets given statements
func (i *Interpreter) Interpret(statements []Stmt) (string, error) {
	var s string
	var err error
	for _, statement := range statements {
		v, er := i.execute(statement)
		err = er
		s = stringfy(v)
		if err != nil {
			i.Runtime.RuntimeError(err)
		}
	}

	return s, err
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
	case GreaterTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqualTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LessTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqualTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BangEqualTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) != right.(float64), nil
	case EqualEqualTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) == right.(float64), nil
	case MinusTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case PlusTT:
		if isFloat64(left) && isFloat64(right) {
			return left.(float64) + right.(float64), nil
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		}

		return nil, RuntimeError.New(expr.Operator, "Operands must be two numbers or two strings.")
	case SlashTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case StarTT:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	// Unreachable.
	return nil, RuntimeError.New(nil, "Unreachable")
}

func (i *Interpreter) visitCallExpr(expr *Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, 0)
	for _, argument := range expr.Arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, RuntimeError.New(expr.Paren, "Can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		return nil, RuntimeError.New(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) visitGetExpr(expr *Get) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}
	if _, ok := object.(*LoxInstance); ok {
		return object.(*LoxInstance).Get(expr.Name)
	}

	return nil, RuntimeError.New(expr.Name, "Only instances have properties.")
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) visitLogicalExpr(expr *Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OrTT {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) visitSetExpr(expr *Set) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if _, ok := object.(*LoxInstance); !ok {
		return nil, RuntimeError.New(expr.Name, "Only instances have fields.")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	li := object.(*LoxInstance)
	li.Set(expr.Name, value)

	return value, nil
}

func (i *Interpreter) visitThisExpr(expr *This) (interface{}, error) {
	return i.lookUpVariable(expr.Keyword, expr)
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
	case BangTT:
		return !i.isTruthy(right), nil
	case MinusTT:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}

	// Unreachable
	return nil, RuntimeError.New(nil, "Unreachable")
}

func (i *Interpreter) visitVariableExpr(expr *Variable) (interface{}, error) {
	// return i.Runtime.Environment.Get(expr.Name)
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name *Token, expr Expr) (interface{}, error) {
	if distance, ok := i.Runtime.Locals[expr]; ok {
		return i.Runtime.Environment.GetAt(distance, name.Lexeme)
	}
	return i.Runtime.Globals.Get(name)
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

// Resolve resolves an expression
func (i *Interpreter) Resolve(expr Expr, depth int) error {
	i.Runtime.Locals[expr] = depth

	return nil
}

func (i *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) (interface{}, error) {
	previous := i.Runtime.Environment
	defer func() { i.Runtime.Environment = previous }()
	i.Runtime.Environment = environment
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) visitBlockStmt(stmt *Block) (interface{}, error) {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.Runtime.Environment))
}

func (i *Interpreter) visitClassStmt(stmt *Class) (interface{}, error) {
	i.Runtime.Environment.Define(stmt.Name.Lexeme, nil)

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := NewLoxFunction(method, i.Runtime.Environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	klass := NewLoxClass(stmt.Name.Lexeme, methods)
	i.Runtime.Environment.Assign(stmt.Name, klass)

	return nil, nil
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

func (i *Interpreter) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
	// return nil, nil
}

func (i *Interpreter) visitFunctionStmt(stmt *Function) (interface{}, error) {
	function := NewLoxFunction(stmt, i.Runtime.Environment, false)
	i.Runtime.Environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) visitIfStmt(stmt *If) (interface{}, error) {
	val, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(val) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) visitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(stringfy(value))
	return nil, nil
}

func (i *Interpreter) visitReturnStmt(stmt *Return) (interface{}, error) {
	var value interface{} = nil
	if stmt.Value != nil {
		var err error
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, NewReturnValue(value)
}

func (i *Interpreter) visitWhileStmt(stmt *While) (interface{}, error) {
	v, _ := i.evaluate(stmt.Condition)
	for ; i.isTruthy(v); v, _ = i.evaluate(stmt.Condition) {
		i.execute(stmt.Body)
	}

	return nil, nil
}

func (i *Interpreter) visitVarStmt(stmt *Var) (interface{}, error) {
	var value interface{} = nil
	if stmt.Initializer != nil {
		v, err := i.evaluate(stmt.Initializer)
		value = v
		if err != nil {
			return nil, err
		}
	}

	i.Runtime.Environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) visitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if distance, ok := i.Runtime.Locals[expr]; ok {
		i.Runtime.Environment.AssignAt(distance, expr.Name, value)
	} else {
		err := i.Runtime.Globals.Assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// ReturnValue is struct of return value
type ReturnValue struct {
	Value interface{}
}

// Error satisfies error interface
func (r *ReturnValue) Error() string {
	return "Return Value error"
}

// NewReturnValue is constructor of ReturValue
func NewReturnValue(value interface{}) *ReturnValue {
	return &ReturnValue{Value: value}
}

func isType(v interface{}, kind reflect.Kind) bool {
	return reflect.ValueOf(v).Kind() == kind
}

func isFloat64(v interface{}) bool {
	return isType(v, reflect.Float64)
}

func isString(v interface{}) bool {
	return isType(v, reflect.String)
}

func stringfy(object interface{}) string {
	if object == nil {
		return "nil"
	}

	return fmt.Sprint(object)
}
