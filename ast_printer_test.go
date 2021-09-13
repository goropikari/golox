package golox_test

import (
	"testing"

	"github.com/goropikari/golox"
	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	ast := golox.NewAstPrinter()

	var tests = []struct {
		name     string
		expected string
		given    []golox.Stmt
	}{
		{
			name:     "binary: -123 * (45.67)",
			expected: "(* (- 123) (group 45.67))",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(
					golox.NewUnary(
						golox.NewToken(golox.MinusTT, "-", nil, 1), golox.NewLiteral(123),
					),
					golox.NewToken(golox.StarTT, "*", nil, 1),
					golox.NewGrouping(golox.NewLiteral(45.67))),
				),
			},
		},
		{
			name:     "logical: -123 and 45.67",
			expected: "(and (- 123) (group 45.67))",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewLogical(
					golox.NewUnary(
						golox.NewToken(golox.MinusTT, "-", nil, 1), golox.NewLiteral(123),
					),
					golox.NewToken(golox.AndTT, "and", nil, 1),
					golox.NewGrouping(golox.NewLiteral(45.67))),
				),
			},
		},
		{
			name:     "assign: x = 123",
			expected: "(assign x 123)",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewAssign(
					golox.NewToken(golox.IdentifierTT, "x", nil, 1),
					golox.NewLiteral(123)),
				),
			},
		},
		{
			name:     "visit variable: x",
			expected: "(variable x)",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewVariable(
					golox.NewToken(golox.IdentifierTT, "x", nil, 1),
				)),
			},
		},
		{
			name:     "if statement",
			expected: "(if (cond true) (thenBranch (variable x)) (elseBranch (variable y)))",
			given: []golox.Stmt{
				golox.NewIf(
					golox.NewLiteral(true),
					golox.NewExpression(
						golox.NewVariable(
							golox.NewToken(golox.IdentifierTT, "x", nil, 1),
						),
					),
					golox.NewExpression(
						golox.NewVariable(
							golox.NewToken(golox.IdentifierTT, "y", nil, 1),
						),
					),
				),
			},
		},
		{
			name:     "while statement",
			expected: "(while (cond 123) (body (print 123)))",
			given: []golox.Stmt{
				golox.NewWhile(
					golox.NewLiteral(123),
					golox.NewPrint(golox.NewLiteral(123)),
				),
			},
		},
		{
			name:     "declare variable: var x = 123",
			expected: "(declare x (initializer 123))",
			given: []golox.Stmt{
				golox.NewVar(
					golox.NewToken(golox.IdentifierTT, "x", nil, 1),
					golox.NewLiteral(123)),
			},
		},
		{
			name:     "block statement: { 123; 987; }",
			expected: "(block (block body 123) (block body 987))",
			given: []golox.Stmt{
				golox.NewBlock(
					[]golox.Stmt{
						golox.NewExpression(
							golox.NewLiteral(123)),
						golox.NewExpression(
							golox.NewLiteral(987)),
					},
				),
			},
		},
		{
			name:     "function",
			expected: "(function f (args (x, y)) (body (1) (2)))",
			given: []golox.Stmt{
				golox.NewFunction(
					golox.NewToken(golox.IdentifierTT, "f", nil, 1),
					[]*golox.Token{
						golox.NewToken(golox.IdentifierTT, "x", nil, 1),
						golox.NewToken(golox.IdentifierTT, "y", nil, 1),
					},
					[]golox.Stmt{
						golox.NewExpression(golox.NewLiteral(1)),
						golox.NewExpression(golox.NewLiteral(2)),
					},
				),
			},
		},
		{
			name:     "class",
			expected: "(class Hoge (function init (args (x)) (body ((set (object (this))(name x)(value (variable x)))))))",
			given: []golox.Stmt{
				// class Hoge:
				//   init(x):
				//     this.x = x
				golox.NewClass(
					golox.NewToken(golox.IdentifierTT, "Hoge", nil, 1),
					nil,
					[]*golox.Function{
						golox.NewFunction(
							golox.NewToken(golox.IdentifierTT, "init", nil, 2),
							[]*golox.Token{
								golox.NewToken(golox.IdentifierTT, "x", nil, 2),
							},
							[]golox.Stmt{
								golox.NewExpression(
									golox.NewSet(
										golox.NewThis(
											golox.NewToken(golox.ThisTT, "this", nil, 3),
										),
										golox.NewToken(golox.IdentifierTT, "x", nil, 3),
										golox.NewVariable(
											golox.NewToken(golox.IdentifierTT, "x", nil, 3),
										),
									),
								),
							},
						).(*golox.Function),
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
