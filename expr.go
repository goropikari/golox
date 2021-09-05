package tlps

type Expr interface {
	Accept(VisitorExpr) (interface{}, error)
	IsType(interface{}) bool
}

type VisitorExpr interface {
	visitAssignExpr(*Assign) (interface{}, error)
	visitBinaryExpr(*Binary) (interface{}, error)
	visitGroupingExpr(*Grouping) (interface{}, error)
	visitLiteralExpr(*Literal) (interface{}, error)
	visitUnaryExpr(*Unary) (interface{}, error)
	visitVariableExpr(*Variable) (interface{}, error)
}

type Assign struct {
	Name  *Token
	Value Expr
}

func NewAssign(name *Token, value Expr) Expr {
	return &Assign{name, value}
}

func (a *Assign) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitAssignExpr(a)
}

func (rec *Assign) IsType(v interface{}) bool {
	switch v.(type) {
	case *Assign:
		return true
	}
	return false
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) Expr {
	return &Binary{left, operator, right}
}

func (b *Binary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitBinaryExpr(b)
}

func (rec *Binary) IsType(v interface{}) bool {
	switch v.(type) {
	case *Binary:
		return true
	}
	return false
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{expression}
}

func (g *Grouping) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitGroupingExpr(g)
}

func (rec *Grouping) IsType(v interface{}) bool {
	switch v.(type) {
	case *Grouping:
		return true
	}
	return false
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Expr {
	return &Literal{value}
}

func (l *Literal) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitLiteralExpr(l)
}

func (rec *Literal) IsType(v interface{}) bool {
	switch v.(type) {
	case *Literal:
		return true
	}
	return false
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(operator *Token, right Expr) Expr {
	return &Unary{operator, right}
}

func (u *Unary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitUnaryExpr(u)
}

func (rec *Unary) IsType(v interface{}) bool {
	switch v.(type) {
	case *Unary:
		return true
	}
	return false
}

type Variable struct {
	Name *Token
}

func NewVariable(name *Token) Expr {
	return &Variable{name}
}

func (v *Variable) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitVariableExpr(v)
}

func (rec *Variable) IsType(v interface{}) bool {
	switch v.(type) {
	case *Variable:
		return true
	}
	return false
}
