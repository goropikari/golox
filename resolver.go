package tlps

import (
	"errors"
)

// Resolver is struct of resolver
type Resolver struct {
	runtime         *Runtime
	Interpreter     *Interpreter
	currentFunction FunctionType
	currentClass    ClassType
}

// FunctionType is current scope function type
type FunctionType int

const (
	NoneFT FunctionType = iota
	FunctionFT
	InitializerFT
	MethodFT
)

type ClassType int

const (
	NoneCT ClassType = iota
	ClassCT
)

// NewResolver is constructor of Resolver
func NewResolver(runtime *Runtime, interpreter *Interpreter) *Resolver {
	return &Resolver{
		runtime:         runtime,
		Interpreter:     interpreter,
		currentFunction: NoneFT,
		currentClass:    NoneCT,
	}
}

func (r *Resolver) visitBlockStmt(stmt *Block) (interface{}, error) {
	if stmt.Typ == NoneBlock {
		r.runtime.ErrorTokenMessage(stmt.Keyword, "unexpected indent")
	}
	r.beginScope()
	r.ResolveStmts(stmt.Statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) visitClassStmt(stmt *Class) (interface{}, error) {
	enclosigClass := r.currentClass
	r.currentClass = ClassCT

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	r.runtime.Scopes.Peek()["this"] = true

	for _, method := range stmt.Methods {
		declaration := MethodFT
		if method.Name.Lexeme == "init" {
			declaration = InitializerFT
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()

	r.currentClass = enclosigClass
	return nil, nil
}

func (r *Resolver) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) visitFunctionStmt(stmt *Function) (interface{}, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionFT)
	return nil, nil
}

func (r *Resolver) visitIfStmt(stmt *If) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		_, err := r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitReturnStmt(stmt *Return) (interface{}, error) {
	if r.currentFunction == NoneFT {
		r.runtime.ErrorTokenMessage(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		if r.currentFunction == InitializerFT {
			r.runtime.ErrorTokenMessage(stmt.Keyword, "Can't return a value from an initializer.")
		}

		_, err := r.resolveExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt *While) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveStmt(stmt.Body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) visitVarStmt(stmt *Var) (interface{}, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(expr *Assign) (interface{}, error) {
	r.resolveExpr(expr.Value)
	return nil, r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) visitBinaryExpr(expr *Binary) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) visitCallExpr(expr *Call) (interface{}, error) {
	_, err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		_, err = r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitGetExpr(expr *Get) (interface{}, error) {
	_, err := r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	_, err := r.resolveExpr(expr.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr *Logical) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitSetExpr(expr *Set) (interface{}, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitThisExpr(expr *This) (interface{}, error) {
	if r.currentClass == NoneCT {
		r.runtime.ErrorTokenMessage(expr.Keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(expr *Unary) (interface{}, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) visitVariableExpr(expr *Variable) (interface{}, error) {
	if !r.runtime.Scopes.IsEmpty() {
		if v, ok := r.runtime.Scopes.Peek()[expr.Name.Lexeme]; !v && ok { // declare variable && not define
			r.runtime.ErrorTokenMessage(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

// ResolveStmts resolves statements
func (r *Resolver) ResolveStmts(stmts []Stmt) (interface{}, error) {
	for _, stmt := range stmts {
		_, err := r.resolveStmt(stmt)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) resolveStmt(stmt Stmt) (interface{}, error) {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) (interface{}, error) {
	return expr.Accept(r)
}

func (r *Resolver) resolveFunction(function *Function, typ FunctionType) (interface{}, error) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typ
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	_, err := r.ResolveStmts(function.Body)
	if err != nil {
		return nil, err
	}
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil, nil
}

func (r *Resolver) beginScope() {
	r.runtime.Scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.runtime.Scopes.Pop()
}

func (r *Resolver) declare(name *Token) {
	if r.runtime.Scopes.IsEmpty() {
		return
	}

	scope := r.runtime.Scopes.Peek()
	if _, ok := scope[name.Lexeme]; ok {
		r.runtime.ErrorTokenMessage(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false

	return
}

func (r *Resolver) define(name *Token) {
	if r.runtime.Scopes.IsEmpty() {
		return
	}
	r.runtime.Scopes.Peek()[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) error {
	for i := 0; i < r.runtime.Scopes.Size(); i++ {
		scope, err := r.runtime.Scopes.Get(i)
		if err != nil {
			return err
		}
		if _, ok := scope[name.Lexeme]; ok {
			r.Interpreter.Resolve(expr, i)
			return nil
		}
	}

	return errors.New("no variable")
}
