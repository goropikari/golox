package mylang

type Expr interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	visitBinaryExpr(*Binary) interface{}
	visitGroupingExpr(*Grouping) interface{}
	visitLiteralExpr(*Literal) interface{}
	visitUnaryExpr(*Unary) interface{}
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(Left Expr, Operator *Token, Right Expr) Expr {
	return &Binary{Left, Operator, Right}
}

func (b *Binary) Accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) Expr {
	return &Grouping{Expression}
}

func (g *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) Expr {
	return &Literal{Value}
}

func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(Operator *Token, Right Expr) Expr {
	return &Unary{Operator, Right}
}

func (u *Unary) Accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(u)
}
