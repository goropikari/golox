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
		"Binary : Left Expr, Operator *Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal : Value interface{}",
		"Unary : Operator *Token, Right Expr",
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
	writer.WriteString("Accept(Visitor) interface{}\n")
	writer.WriteString("}\n\n")

	defineVisitor(writer, baseName, types)

	for _, typ := range types {
		t := strings.Split(typ, ":")
		className := strings.TrimSpace(t[0])
		fields := strings.TrimSpace(t[1])
		defineType(writer, baseName, className, fields)
	}

	return nil
}

func defineType(writer *bufio.Writer, baseName, className, fields string) {
	writer.WriteString("type " + className + " struct {\n")
	fieldList := strings.Split(fields, ", ")
	for _, field := range fieldList {
		writer.WriteString(field + "\n")
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

	writer.WriteString("func (" + strings.ToLower(string(className[0])) + " *" + className + ") Accept(visitor Visitor) interface{} {\n")
	writer.WriteString("return visitor.visit" + className + baseName + "(" + strings.ToLower(string(className[0])) + ")")
	writer.WriteString("}\n")

}

func defineVisitor(writer *bufio.Writer, baseName string, types []string) {
	writer.WriteString("type Visitor interface {\n")
	for _, typ := range types {
		typName := strings.TrimSpace(strings.Split(typ, ":")[0])
		writer.WriteString("visit" + typName + baseName + "(*" + typName + ")" + " interface{}\n")
	}
	writer.WriteString("}\n\n")
}
