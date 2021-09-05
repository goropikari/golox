package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: generate_ast <output directory>\n")
		os.Exit(64)
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Assign : name *Token, value Expr",
		"Binary : left Expr, operator *Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value interface{}",
		"Logical : left Expr, operator *Token, right Expr",
		"Unary : operator *Token, right Expr",
		"Variable : name *Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block : statements []Stmt",
		"Expression: expression Expr",
		"If_ : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print_ : expression Expr",
		"Var_ : name *Token, initializer Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.WriteString("package tlps\n")

	writer.WriteString("type " + baseName + " interface {")
	writer.WriteString("Accept(Visitor" + baseName + ") (interface{}, error)\n")
	writer.WriteString("IsType(interface{}) bool\n")
	writer.WriteString("}\n\n")

	defineVisitor(writer, baseName, types)

	for _, typ := range types {
		t := strings.Split(typ, ":")
		className := strings.TrimSpace(t[0])
		fields := strings.TrimSpace(t[1])
		defineType(writer, baseName, className, fields)
		defineIsType(writer, className)
	}

	return nil
}

func defineType(writer *bufio.Writer, baseName, className, fields string) {
	writer.WriteString("type " + className + " struct {\n")
	fieldList := strings.Split(fields, ", ")
	for _, field := range fieldList {
		vs := strings.Split(field, " ")
		for i, v := range vs {
			if i == 0 {
				writer.WriteString(strings.Title(v) + " ")
			} else {
				writer.WriteString(v + "\n")
			}
		}
	}
	writer.WriteString("}\n\n")

	writer.WriteString("func New" + className + "(" + fields + ") " + baseName + "{\n")
	writer.WriteString("return &" + className + "{")
	args := make([]string, 0)
	for _, field := range fieldList {
		name := strings.Split(field, " ")[0]
		args = append(args, name)
	}
	writer.WriteString(strings.Join(args, ","))

	writer.WriteString("}\n")
	writer.WriteString("}\n\n")

	writer.WriteString("func (" + strings.ToLower(string(className[0])) + " *" + className + ") Accept(visitor Visitor" + baseName + ") (interface{}, error) {\n")
	writer.WriteString("return visitor.visit" + className + baseName + "(" + strings.ToLower(string(className[0])) + ")\n")
	writer.WriteString("}\n\n")
}

func defineVisitor(writer *bufio.Writer, baseName string, types []string) {
	writer.WriteString("type Visitor" + baseName + " interface {\n")
	for _, typ := range types {
		typName := strings.TrimSpace(strings.Split(typ, ":")[0])
		writer.WriteString("visit" + typName + baseName + "(*" + typName + ")" + " (interface{}, error)\n")
	}
	writer.WriteString("}\n\n")
}

func defineIsType(writer *bufio.Writer, className string) {
	writer.WriteString("func (rec *" + className + ") IsType(v interface{}) bool {\n")
	writer.WriteString("switch v.(type) {\n")
	writer.WriteString("case *" + className + ":\n")
	writer.WriteString("return true\n")
	writer.WriteString("}\n")
	writer.WriteString("return false")
	writer.WriteString("}\n")
}
