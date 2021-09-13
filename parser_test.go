package golox_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/golox"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	r := golox.NewRuntime()

	var tests = []struct {
		name     string
		expected []golox.Stmt
		given    golox.TokenList
	}{
		{
			name: "1 + 2 * 3",
			expected: []golox.Stmt{
				golox.NewExpression(
					golox.NewBinary(
						golox.NewLiteral(1.0),
						golox.NewToken(golox.PlusTT, "+", nil, 1),
						golox.NewBinary(
							golox.NewLiteral(2.0),
							golox.NewToken(golox.StarTT, "*", nil, 1),
							golox.NewLiteral(3.0),
						),
					),
				),
			},
			given: golox.TokenList{
				golox.NewToken(golox.NumberTT, "1", 1.0, 1),
				golox.NewToken(golox.PlusTT, "+", nil, 1),
				golox.NewToken(golox.NumberTT, "2", 2.0, 1),
				golox.NewToken(golox.StarTT, "*", nil, 1),
				golox.NewToken(golox.NumberTT, "3", 3.0, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
		},
		{
			name: "if (true) { print 1; }",
			expected: []golox.Stmt{
				golox.NewIf(
					golox.NewLiteral(true),
					golox.NewBlock(
						[]golox.Stmt{
							golox.NewPrint(golox.NewLiteral(1.0)),
						},
					),
					nil,
				),
			},
			given: golox.TokenList{
				golox.NewToken(golox.IfTT, "if", nil, 1),
				golox.NewToken(golox.LeftParenTT, "(", nil, 1),
				golox.NewToken(golox.TrueTT, "true", nil, 1),
				golox.NewToken(golox.RightParenTT, ")", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 2),
				golox.NewToken(golox.PrintTT, "print", nil, 2),
				golox.NewToken(golox.NumberTT, "1", 1.0, 2),
				golox.NewToken(golox.SemicolonTT, ";", nil, 2),
				golox.NewToken(golox.RightBraceTT, "}", nil, 2),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
		},
		{
			name: "if true { print 1; } else { print 2; }",
			expected: []golox.Stmt{
				golox.NewIf(
					golox.NewLiteral(true),
					golox.NewBlock(
						[]golox.Stmt{
							golox.NewPrint(golox.NewLiteral(1.0)),
						},
					),
					golox.NewBlock(
						[]golox.Stmt{
							golox.NewPrint(golox.NewLiteral(2.0)),
						},
					),
				),
			},
			given: golox.TokenList{
				// then branch
				golox.NewToken(golox.IfTT, "if", nil, 1),
				golox.NewToken(golox.LeftParenTT, "(", nil, 1),
				golox.NewToken(golox.TrueTT, "true", nil, 1),
				golox.NewToken(golox.RightParenTT, ")", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 2),
				golox.NewToken(golox.PrintTT, "print", nil, 2),
				golox.NewToken(golox.NumberTT, "1", 1.0, 2),
				golox.NewToken(golox.SemicolonTT, ";", nil, 2),
				golox.NewToken(golox.RightBraceTT, "}", nil, 2),

				// else branch
				golox.NewToken(golox.ElseTT, "else", nil, 2),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 4),
				golox.NewToken(golox.PrintTT, "print", nil, 4),
				golox.NewToken(golox.NumberTT, "2", 2.0, 4),
				golox.NewToken(golox.SemicolonTT, ";", nil, 4),
				golox.NewToken(golox.RightBraceTT, "}", nil, 4),
				golox.NewToken(golox.EOFTT, "", nil, 4),
			},
		},
		{
			name: "for (var i = 0; i < 5; i = i + 1) { print i; }",
			expected: []golox.Stmt{
				golox.NewBlock(
					[]golox.Stmt{
						golox.NewVar(
							golox.NewToken(golox.IdentifierTT, "i", nil, 1),
							golox.NewLiteral(0.0),
						),
						golox.NewWhile(
							golox.NewBinary(
								golox.NewVariable(golox.NewToken(golox.IdentifierTT, "i", nil, 1)),
								golox.NewToken(golox.LessTT, "<", nil, 1),
								golox.NewLiteral(5.0),
							),
							golox.NewBlock(
								[]golox.Stmt{
									golox.NewBlock(
										[]golox.Stmt{
											golox.NewPrint(
												golox.NewVariable(
													golox.NewToken(golox.IdentifierTT, "i", nil, 2),
												),
											),
										},
									),
									golox.NewExpression(
										golox.NewAssign(
											golox.NewToken(golox.IdentifierTT, "i", nil, 1),
											golox.NewBinary(
												golox.NewVariable(
													golox.NewToken(golox.IdentifierTT, "i", nil, 1),
												),
												golox.NewToken(golox.PlusTT, "+", nil, 1),
												golox.NewLiteral(1.0),
											),
										),
									),
								},
							),
						),
					},
				),
			},
			given: golox.TokenList{
				golox.NewToken(golox.ForTT, "for", nil, 1),
				golox.NewToken(golox.LeftParenTT, "(", nil, 1),
				golox.NewToken(golox.VarTT, "var", nil, 1),
				golox.NewToken(golox.IdentifierTT, "i", nil, 1),
				golox.NewToken(golox.EqualTT, "=", nil, 1),
				golox.NewToken(golox.NumberTT, "0", 0.0, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.IdentifierTT, "i", nil, 1),
				golox.NewToken(golox.LessTT, "<", nil, 1),
				golox.NewToken(golox.NumberTT, "5", 5.0, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.IdentifierTT, "i", nil, 1),
				golox.NewToken(golox.EqualTT, "=", nil, 1),
				golox.NewToken(golox.IdentifierTT, "i", nil, 1),
				golox.NewToken(golox.PlusTT, "+", nil, 1),
				golox.NewToken(golox.NumberTT, "1", 1.0, 1),
				golox.NewToken(golox.RightParenTT, ")", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 2),
				golox.NewToken(golox.PrintTT, "print", nil, 2),
				golox.NewToken(golox.IdentifierTT, "i", nil, 2),
				golox.NewToken(golox.SemicolonTT, ";", nil, 2),
				golox.NewToken(golox.RightBraceTT, "}", nil, 2),
				golox.NewToken(golox.EOFTT, "", nil, 2),
			},
		},
		{
			name: "function",
			expected: []golox.Stmt{
				// fun f(x, y) {
				//   return x + y;
				// }
				golox.NewFunction(
					golox.NewToken(golox.IdentifierTT, "f", nil, 1),
					[]*golox.Token{
						golox.NewToken(golox.IdentifierTT, "x", nil, 1),
						golox.NewToken(golox.IdentifierTT, "y", nil, 1),
					},
					[]golox.Stmt{
						golox.NewReturn(
							golox.NewToken(golox.ReturnTT, "return", nil, 2),
							golox.NewBinary(
								golox.NewVariable(golox.NewToken(golox.IdentifierTT, "x", nil, 2)),
								golox.NewToken(golox.PlusTT, "+", nil, 2),
								golox.NewVariable(golox.NewToken(golox.IdentifierTT, "y", nil, 2)),
							),
						),
					},
				),
			},
			given: []*golox.Token{
				golox.NewToken(golox.FunTT, "fun", nil, 1),
				golox.NewToken(golox.IdentifierTT, "f", nil, 1),
				golox.NewToken(golox.LeftParenTT, "(", nil, 1),
				golox.NewToken(golox.IdentifierTT, "x", nil, 1),
				golox.NewToken(golox.CommaTT, ",", nil, 1),
				golox.NewToken(golox.IdentifierTT, "y", nil, 1),
				golox.NewToken(golox.RightParenTT, ")", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 2),
				golox.NewToken(golox.ReturnTT, "return", nil, 2),
				golox.NewToken(golox.IdentifierTT, "x", nil, 2),
				golox.NewToken(golox.PlusTT, "+", nil, 2),
				golox.NewToken(golox.IdentifierTT, "y", nil, 2),
				golox.NewToken(golox.SemicolonTT, ";", nil, 2),
				golox.NewToken(golox.RightBraceTT, "}", nil, 2),
				golox.NewToken(golox.EOFTT, "", nil, 2),
			},
		},
		{
			name: "class",
			expected: []golox.Stmt{
				// class Hoge {
				//   init(x) {
				//     this.x = x;
				//   }
				// }
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
			given: []*golox.Token{
				golox.NewToken(golox.ClassTT, "class", nil, 1),
				golox.NewToken(golox.IdentifierTT, "Hoge", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 2),
				golox.NewToken(golox.IdentifierTT, "init", nil, 2),
				golox.NewToken(golox.LeftParenTT, "(", nil, 2),
				golox.NewToken(golox.IdentifierTT, "x", nil, 2),
				golox.NewToken(golox.RightParenTT, ")", nil, 2),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 3),
				golox.NewToken(golox.ThisTT, "this", nil, 3),
				golox.NewToken(golox.DotTT, ".", nil, 3),
				golox.NewToken(golox.IdentifierTT, "x", nil, 3),
				golox.NewToken(golox.EqualTT, "=", nil, 3),
				golox.NewToken(golox.IdentifierTT, "x", nil, 3),
				golox.NewToken(golox.SemicolonTT, ";", nil, 3),
				golox.NewToken(golox.RightBraceTT, "}", nil, 3),
				golox.NewToken(golox.RightBraceTT, "}", nil, 3),
				golox.NewToken(golox.EOFTT, "", nil, 3),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := golox.NewParser(r, tt.given)
			actual, _ := parser.Parse()
			ast := golox.NewAstPrinter()
			fmt.Println(ast.Print(actual))
			fmt.Println(ast.Print(tt.expected))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
