package tlps

import (
	"bytes"
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(stmts []Stmt) (string, error) {
	val, err := stmts[0].Accept(ap)
	return val.(string), err
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesizeExpr("group", expr.Expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (ap *AstPrinter) visitLogicalExpr(expr *Logical) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) visitAssignExpr(expr *Assign) (interface{}, error) {
	return ap.parenthesizeExpr("assign "+expr.Name.Lexeme, expr.Value)
}

func (ap *AstPrinter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return ap.parenthesizeExpr("variable", NewLiteral(expr.Name.Lexeme))
}

func (ap *AstPrinter) visitBlockStmt(b *Block) (interface{}, error) {
	body := make([]string, 0)
	for _, stmt := range b.Statements {
		s, err := ap.parenthesizeStmt("block body", stmt)
		if err != nil {
			return "", err
		}
		body = append(body, s)
	}
	return "(block " + strings.Join(body, " ") + ")", nil
}

func (ap *AstPrinter) visitExpressionStmt(e *Expression) (interface{}, error) {
	return e.Expression.Accept(ap)
}

func (ap *AstPrinter) visitIf_Stmt(i *If_) (interface{}, error) {
	cond, err := ap.parenthesizeExpr("cond", i.Condition)
	if err != nil {
		return "", nil
	}
	thenBranch, err := ap.parenthesizeStmt("thenBranch", i.ThenBranch)
	if err != nil {
		return "", nil
	}
	var elseBranch string
	if i.ElseBranch != nil {
		elseBranch, err = ap.parenthesizeStmt("elseBranch", i.ElseBranch)
		if err != nil {
			return "", err
		}
	}

	return "(if " + cond + " " + thenBranch + " " + elseBranch + ")", nil
}

func (ap *AstPrinter) visitPrint_Stmt(p *Print_) (interface{}, error) {
	return ap.parenthesizeExpr("print", p.Expression)
}

func (ap *AstPrinter) visitWhile_Stmt(p *While_) (interface{}, error) {
	cond, err := ap.parenthesizeExpr("cond", p.Condition)
	if err != nil {
		return "", err
	}
	body, err := ap.parenthesizeStmt("body", p.Body)
	if err != nil {
		return "", nil
	}
	return "(while " + cond + " " + body + ")", nil
}

func (ap *AstPrinter) visitVar_Stmt(v *Var_) (interface{}, error) {
	initializer, err := ap.parenthesizeExpr("init", v.Initializer)
	if err != nil {
		return "", nil
	}
	return "(declare " + v.Name.Lexeme + " " + initializer + ")", nil
}

func (ap *AstPrinter) parenthesizeExpr(name string, exprs ...Expr) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, expr := range exprs {
		buf.WriteString(" ")
		s, _ := expr.Accept(ap)
		buf.WriteString(s.(string))
	}
	buf.WriteString(")")

	return buf.String(), nil
}

func (ap *AstPrinter) parenthesizeStmt(name string, stmts ...Stmt) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, stmt := range stmts {
		buf.WriteString(" ")
		s, _ := stmt.Accept(ap)
		buf.WriteString(s.(string))
	}
	buf.WriteString(")")

	return buf.String(), nil
}
