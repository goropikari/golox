package tlps

type Expr interface {
	Accept(VisitorExpr) (interface{}, error)
	IsType(interface{}) bool
}

type VisitorExpr interface {
	visitAssignExpr(*Assign) (interface{}, error)
	visitBinaryExpr(*Binary) (interface{}, error)
	visitCallExpr(*Call) (interface{}, error)
	visitGetExpr(*Get) (interface{}, error)
	visitGroupingExpr(*Grouping) (interface{}, error)
	visitLiteralExpr(*Literal) (interface{}, error)
	visitLogicalExpr(*Logical) (interface{}, error)
	visitSetExpr(*Set) (interface{}, error)
	visitSuperExpr(*Super) (interface{}, error)
	visitThisExpr(*This) (interface{}, error)
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

type Call struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *Token, arguments []Expr) Expr {
	return &Call{callee, paren, arguments}
}

func (c *Call) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitCallExpr(c)
}

func (rec *Call) IsType(v interface{}) bool {
	switch v.(type) {
	case *Call:
		return true
	}
	return false
}

type Get struct {
	Object Expr
	Name   *Token
}

func NewGet(object Expr, name *Token) Expr {
	return &Get{object, name}
}

func (g *Get) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitGetExpr(g)
}

func (rec *Get) IsType(v interface{}) bool {
	switch v.(type) {
	case *Get:
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

type Logical struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewLogical(left Expr, operator *Token, right Expr) Expr {
	return &Logical{left, operator, right}
}

func (l *Logical) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitLogicalExpr(l)
}

func (rec *Logical) IsType(v interface{}) bool {
	switch v.(type) {
	case *Logical:
		return true
	}
	return false
}

type Set struct {
	Object Expr
	Name   *Token
	Value  Expr
}

func NewSet(object Expr, name *Token, value Expr) Expr {
	return &Set{object, name, value}
}

func (s *Set) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitSetExpr(s)
}

func (rec *Set) IsType(v interface{}) bool {
	switch v.(type) {
	case *Set:
		return true
	}
	return false
}

type Super struct {
	Keyword *Token
	Method  *Token
}

func NewSuper(keyword *Token, method *Token) Expr {
	return &Super{keyword, method}
}

func (s *Super) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitSuperExpr(s)
}

func (rec *Super) IsType(v interface{}) bool {
	switch v.(type) {
	case *Super:
		return true
	}
	return false
}

type This struct {
	Keyword *Token
}

func NewThis(keyword *Token) Expr {
	return &This{keyword}
}

func (t *This) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.visitThisExpr(t)
}

func (rec *This) IsType(v interface{}) bool {
	switch v.(type) {
	case *This:
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
