package tlps

type Stmt interface {
	Accept(VisitorStmt) (interface{}, error)
	IsType(interface{}) bool
}

type VisitorStmt interface {
	visitBlockStmt(*Block) (interface{}, error)
	visitExpressionStmt(*Expression) (interface{}, error)
	visitPrint_Stmt(*Print_) (interface{}, error)
	visitVar_Stmt(*Var_) (interface{}, error)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) Stmt {
	return &Block{statements}
}

func (b *Block) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitBlockStmt(b)
}

func (rec *Block) IsType(v interface{}) bool {
	switch v.(type) {
	case *Block:
		return true
	}
	return false
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{expression}
}

func (e *Expression) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitExpressionStmt(e)
}

func (rec *Expression) IsType(v interface{}) bool {
	switch v.(type) {
	case *Expression:
		return true
	}
	return false
}

type Print_ struct {
	Expression Expr
}

func NewPrint_(expression Expr) Stmt {
	return &Print_{expression}
}

func (p *Print_) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitPrint_Stmt(p)
}

func (rec *Print_) IsType(v interface{}) bool {
	switch v.(type) {
	case *Print_:
		return true
	}
	return false
}

type Var_ struct {
	Name        *Token
	Initializer Expr
}

func NewVar_(name *Token, initializer Expr) Stmt {
	return &Var_{name, initializer}
}

func (v *Var_) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitVar_Stmt(v)
}

func (rec *Var_) IsType(v interface{}) bool {
	switch v.(type) {
	case *Var_:
		return true
	}
	return false
}
