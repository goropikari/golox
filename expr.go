package tlps

type Expr interface {
	Accept(Visitor) (interface{}, error)
}

type Visitor interface {
	visitBinaryExpr(*Binary) (interface{}, error)
	visitGroupingExpr(*Grouping) (interface{}, error)
	visitLiteralExpr(*Literal) (interface{}, error)
	visitUnaryExpr(*Unary) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(Left Expr, Operator *Token, Right Expr) Expr {
	return &Binary{Left, Operator, Right}
}

func (b *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) Expr {
	return &Grouping{Expression}
}

func (g *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) Expr {
	return &Literal{Value}
}

func (l *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitLiteralExpr(l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(Operator *Token, Right Expr) Expr {
	return &Unary{Operator, Right}
}

func (u *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitUnaryExpr(u)
}
