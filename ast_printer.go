package tlps

import (
	"bytes"
	"fmt"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(expr Expr) (string, error) {
	val, err := expr.Accept(ap)
	return val.(string), err
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.Expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
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
