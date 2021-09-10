package tlps_test

import (
	"testing"

	"github.com/goropikari/tlps"
	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	ast := tlps.NewAstPrinter()

	var tests = []struct {
		name     string
		expected string
		given    []tlps.Stmt
	}{
		{
			name:     "binary: -123 * (45.67)",
			expected: "(* (- 123) (group 45.67))",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(
					tlps.NewUnary(
						tlps.NewToken(tlps.MinusTT, "-", nil, 1), tlps.NewLiteral(123),
					),
					tlps.NewToken(tlps.StarTT, "*", nil, 1),
					tlps.NewGrouping(tlps.NewLiteral(45.67))),
				),
			},
		},
		{
			name:     "logical: -123 and 45.67",
			expected: "(and (- 123) (group 45.67))",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewLogical(
					tlps.NewUnary(
						tlps.NewToken(tlps.MinusTT, "-", nil, 1), tlps.NewLiteral(123),
					),
					tlps.NewToken(tlps.AndTT, "and", nil, 1),
					tlps.NewGrouping(tlps.NewLiteral(45.67))),
				),
			},
		},
		{
			name:     "assign: x = 123",
			expected: "(assign x 123)",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewAssign(
					tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
					tlps.NewLiteral(123)),
				),
			},
		},
		{
			name:     "visit variable: x",
			expected: "(variable x)",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewVariable(
					tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
				)),
			},
		},
		{
			name:     "if statement",
			expected: "(if (cond true) (thenBranch (variable x)) (elseBranch (variable y)))",
			given: []tlps.Stmt{
				tlps.NewIf(
					tlps.NewLiteral(true),
					tlps.NewExpression(
						tlps.NewVariable(
							tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
						),
					),
					tlps.NewExpression(
						tlps.NewVariable(
							tlps.NewToken(tlps.IdentifierTT, "y", nil, 1),
						),
					),
				),
			},
		},
		{
			name:     "print 123",
			expected: "(print 123)",
			given: []tlps.Stmt{
				tlps.NewPrint(tlps.NewLiteral(123)),
			},
		},
		{
			name:     "while statement",
			expected: "(while (cond 123) (body (print 123)))",
			given: []tlps.Stmt{
				tlps.NewWhile(
					tlps.NewLiteral(123),
					tlps.NewPrint(tlps.NewLiteral(123)),
				),
			},
		},
		{
			name:     "declare variable: var x = 123",
			expected: "(declare x (init 123))",
			given: []tlps.Stmt{
				tlps.NewVar(
					tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
					tlps.NewLiteral(123)),
			},
		},
		{
			name:     "block statement: { 123; 987; }",
			expected: "(block (block body 123) (block body 987))",
			given: []tlps.Stmt{
				tlps.NewBlock(
					[]tlps.Stmt{
						tlps.NewExpression(
							tlps.NewLiteral(123)),
						tlps.NewExpression(
							tlps.NewLiteral(987)),
					},
				),
			},
		},
		{
			name:     "function",
			expected: "(function (args (x, y)) (body (1) (2)))",
			given: []tlps.Stmt{
				tlps.NewFunction(
					tlps.NewToken(tlps.IdentifierTT, "f", nil, 1),
					[]*tlps.Token{
						tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
						tlps.NewToken(tlps.IdentifierTT, "y", nil, 1),
					},
					[]tlps.Stmt{
						tlps.NewExpression(tlps.NewLiteral(1)),
						tlps.NewExpression(tlps.NewLiteral(2)),
					},
				),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := ast.Print(tt.given)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
