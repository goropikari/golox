package mylang

import (
	"bytes"
	"fmt"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary) interface{} {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping) interface{} {
	return ap.parenthesize("group", expr.Expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary) interface{} {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, expr := range exprs {
		buf.WriteString(" ")
		buf.WriteString(expr.Accept(ap).(string))
	}
	buf.WriteString(")")

	return buf.String()
}
