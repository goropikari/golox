package tlps_test

import (
	"testing"

	"github.com/goropikari/tlps"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	r := tlps.NewRuntime()

	var tests = []struct {
		name     string
		expected []tlps.Stmt
		given    tlps.TokenList
	}{
		{
			name: "1 + 2 * 3",
			expected: []tlps.Stmt{
				tlps.NewExpression(
					tlps.NewBinary(
						tlps.NewLiteral(1.0),
						tlps.NewToken(tlps.PlusTT, "+", nil, 1),
						tlps.NewBinary(
							tlps.NewLiteral(2.0),
							tlps.NewToken(tlps.StarTT, "*", nil, 1),
							tlps.NewLiteral(3.0),
						),
					),
				),
			},
			given: tlps.TokenList{
				tlps.NewToken(tlps.NumberTT, "1", 1.0, 1),
				tlps.NewToken(tlps.PlusTT, "+", nil, 1),
				tlps.NewToken(tlps.NumberTT, "2", 2.0, 1),
				tlps.NewToken(tlps.StarTT, "*", nil, 1),
				tlps.NewToken(tlps.NumberTT, "3", 3.0, 1),
				tlps.NewToken(tlps.SemicolonTT, ";", nil, 1),
				tlps.NewToken(tlps.EOFTT, "", nil, 1),
			},
		},
		{
			name: "if true:\n  print 1",
			expected: []tlps.Stmt{
				tlps.NewIf_(
					tlps.NewLiteral(true),
					tlps.NewBlock([]tlps.Stmt{
						tlps.NewPrint_(tlps.NewLiteral(1.0)),
					}),
					nil,
				),
			},
			given: tlps.TokenList{
				tlps.NewToken(tlps.IfTT, "if", nil, 1),
				tlps.NewToken(tlps.TrueTT, "true", nil, 1),
				tlps.NewToken(tlps.ColonTT, ":", nil, 1),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 2),
				tlps.NewToken(tlps.PrintTT, "print", nil, 2),
				tlps.NewToken(tlps.NumberTT, "1", 1.0, 2),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 1),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 2),
				tlps.NewToken(tlps.EOFTT, "", nil, 1),
			},
		},
		{
			name: "if true:\n  print 1\n  else:  print 2\n",
			expected: []tlps.Stmt{
				tlps.NewIf_(
					tlps.NewLiteral(true),
					tlps.NewBlock([]tlps.Stmt{
						tlps.NewPrint_(tlps.NewLiteral(1.0)),
					}),
					tlps.NewBlock([]tlps.Stmt{
						tlps.NewPrint_(tlps.NewLiteral(2.0)),
					}),
				),
			},
			given: tlps.TokenList{
				// then branch
				tlps.NewToken(tlps.IfTT, "if", nil, 1),
				tlps.NewToken(tlps.TrueTT, "true", nil, 1),
				tlps.NewToken(tlps.ColonTT, ":", nil, 1),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 2),
				tlps.NewToken(tlps.PrintTT, "print", nil, 2),
				tlps.NewToken(tlps.NumberTT, "1", 1.0, 2),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 1),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 2),

				// else branch
				tlps.NewToken(tlps.ElseTT, "else", nil, 2),
				tlps.NewToken(tlps.ColonTT, ":", nil, 2),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 2),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 3),
				tlps.NewToken(tlps.PrintTT, "print", nil, 3),
				tlps.NewToken(tlps.NumberTT, "2", 2.0, 3),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 3),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 3),
				tlps.NewToken(tlps.EOFTT, "", nil, 3),
			},
		},
		{
			name: "for var i = 0; i < 5; i = i + 1:\n  print i\n",
			expected: []tlps.Stmt{
				tlps.NewBlock(
					[]tlps.Stmt{
						tlps.NewVar_(
							tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
							tlps.NewLiteral(0.0),
						),
						tlps.NewWhile_(
							tlps.NewBinary(
								tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "i", nil, 1)),
								tlps.NewToken(tlps.LessTT, "<", nil, 1),
								tlps.NewLiteral(5.0),
							),
							tlps.NewBlock([]tlps.Stmt{
								tlps.NewBlock([]tlps.Stmt{
									tlps.NewPrint_(
										tlps.NewVariable(
											tlps.NewToken(tlps.IdentifierTT, "i", nil, 2),
										),
									),
								}),
								tlps.NewExpression(
									tlps.NewAssign(
										tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
										tlps.NewBinary(
											tlps.NewVariable(
												tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
											),
											tlps.NewToken(tlps.PlusTT, "+", nil, 1),
											tlps.NewLiteral(1.0),
										),
									),
								),
							}),
						),
					},
				),
			},
			given: tlps.TokenList{
				tlps.NewToken(tlps.ForTT, "for", nil, 1),
				tlps.NewToken(tlps.VarTT, "var", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
				tlps.NewToken(tlps.EqualTT, "=", nil, 1),
				tlps.NewToken(tlps.NumberTT, "0", 0.0, 1),
				tlps.NewToken(tlps.SemicolonTT, ";", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
				tlps.NewToken(tlps.LessTT, "<", nil, 1),
				tlps.NewToken(tlps.NumberTT, "5", 5.0, 1),
				tlps.NewToken(tlps.SemicolonTT, ";", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
				tlps.NewToken(tlps.EqualTT, "=", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "i", nil, 1),
				tlps.NewToken(tlps.PlusTT, "+", nil, 1),
				tlps.NewToken(tlps.NumberTT, "1", 1.0, 1),
				tlps.NewToken(tlps.ColonTT, ":", nil, 1),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 2),
				tlps.NewToken(tlps.PrintTT, "print", nil, 2),
				tlps.NewToken(tlps.IdentifierTT, "i", nil, 2),
				tlps.NewToken(tlps.NewlineTT, "\n", nil, 2),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 2),
				tlps.NewToken(tlps.EOFTT, "", nil, 2),
			},
		},
		{
			name: "function",
			expected: []tlps.Stmt{
				// fun f(x, y):
				//   return x + y
				tlps.NewFunction(
					tlps.NewToken(tlps.IdentifierTT, "f", nil, 1),
					[]*tlps.Token{
						tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
						tlps.NewToken(tlps.IdentifierTT, "y", nil, 1),
					},
					[]tlps.Stmt{
						tlps.NewReturn_(
							tlps.NewToken(tlps.ReturnTT, "return", nil, 2),
							tlps.NewBinary(
								tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "x", nil, 2)),
								tlps.NewToken(tlps.PlusTT, "+", nil, 2),
								tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "y", nil, 2)),
							),
						),
					},
				),
			},
			given: []*tlps.Token{
				tlps.NewToken(tlps.FunTT, "fun", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "f", nil, 1),
				tlps.NewToken(tlps.LeftParenTT, "(", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
				tlps.NewToken(tlps.CommaTT, ",", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "y", nil, 1),
				tlps.NewToken(tlps.RightParenTT, ")", nil, 1),
				tlps.NewToken(tlps.ColonTT, ":", nil, 1),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 1),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 2),
				tlps.NewToken(tlps.ReturnTT, "return", nil, 2),
				tlps.NewToken(tlps.IdentifierTT, "x", nil, 2),
				tlps.NewToken(tlps.PlusTT, "+", nil, 2),
				tlps.NewToken(tlps.IdentifierTT, "y", nil, 2),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 2),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 2),
				tlps.NewToken(tlps.EOFTT, "", nil, 2),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := tlps.NewParser(r, tt.given)
			actual, _ := parser.Parse()
			// ast := tlps.NewAstPrinter()
			// fmt.Println(ast.Print(actual))
			// fmt.Println(ast.Print(tt.expected))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
