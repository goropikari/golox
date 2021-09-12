package tlps

type Stmt interface {
	Accept(VisitorStmt) (interface{}, error)
	IsType(interface{}) bool
}

type VisitorStmt interface {
	visitBlockStmt(*Block) (interface{}, error)
	visitClassStmt(*Class) (interface{}, error)
	visitExpressionStmt(*Expression) (interface{}, error)
	visitFunctionStmt(*Function) (interface{}, error)
	visitIfStmt(*If) (interface{}, error)
	visitReturnStmt(*Return) (interface{}, error)
	visitVarStmt(*Var) (interface{}, error)
	visitWhileStmt(*While) (interface{}, error)
}

type Block struct {
	Statements []Stmt
	Keyword    *Token
	Typ        BlockType
}

func NewBlock(statements []Stmt, keyword *Token, typ BlockType) Stmt {
	return &Block{statements, keyword, typ}
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

type Class struct {
	Name    *Token
	Methods []*Function
}

func NewClass(name *Token, methods []*Function) Stmt {
	return &Class{name, methods}
}

func (c *Class) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitClassStmt(c)
}

func (rec *Class) IsType(v interface{}) bool {
	switch v.(type) {
	case *Class:
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

type Function struct {
	Name   *Token
	Params []*Token
	Body   []Stmt
}

func NewFunction(name *Token, params []*Token, body []Stmt) Stmt {
	return &Function{name, params, body}
}

func (f *Function) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitFunctionStmt(f)
}

func (rec *Function) IsType(v interface{}) bool {
	switch v.(type) {
	case *Function:
		return true
	}
	return false
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return &If{condition, thenBranch, elseBranch}
}

func (i *If) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitIfStmt(i)
}

func (rec *If) IsType(v interface{}) bool {
	switch v.(type) {
	case *If:
		return true
	}
	return false
}

type Return struct {
	Keyword *Token
	Value   Expr
}

func NewReturn(keyword *Token, value Expr) Stmt {
	return &Return{keyword, value}
}

func (r *Return) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitReturnStmt(r)
}

func (rec *Return) IsType(v interface{}) bool {
	switch v.(type) {
	case *Return:
		return true
	}
	return false
}

type Var struct {
	Name        *Token
	Initializer Expr
}

func NewVar(name *Token, initializer Expr) Stmt {
	return &Var{name, initializer}
}

func (v *Var) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitVarStmt(v)
}

func (rec *Var) IsType(v interface{}) bool {
	switch v.(type) {
	case *Var:
		return true
	}
	return false
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return &While{condition, body}
}

func (w *While) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.visitWhileStmt(w)
}

func (rec *While) IsType(v interface{}) bool {
	switch v.(type) {
	case *While:
		return true
	}
	return false
}
