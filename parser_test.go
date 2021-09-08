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
						tlps.NewToken(tlps.Plus, "+", nil, 1),
						tlps.NewBinary(
							tlps.NewLiteral(2.0),
							tlps.NewToken(tlps.Star, "*", nil, 1),
							tlps.NewLiteral(3.0),
						),
					),
				),
			},
			given: tlps.TokenList{
				tlps.NewToken(tlps.Number, "1", 1.0, 1),
				tlps.NewToken(tlps.Plus, "+", nil, 1),
				tlps.NewToken(tlps.Number, "2", 2.0, 1),
				tlps.NewToken(tlps.Star, "*", nil, 1),
				tlps.NewToken(tlps.Number, "3", 3.0, 1),
				tlps.NewToken(tlps.Semicolon, ";", nil, 1),
				tlps.NewToken(tlps.EOF, "", nil, 1),
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
				tlps.NewToken(tlps.If, "if", nil, 1),
				tlps.NewToken(tlps.True, "true", nil, 1),
				tlps.NewToken(tlps.Colon, ":", nil, 1),
				tlps.NewToken(tlps.Newline, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 2),
				tlps.NewToken(tlps.Print, "print", nil, 2),
				tlps.NewToken(tlps.Number, "1", 1.0, 2),
				tlps.NewToken(tlps.Newline, "\n", nil, 1),
				tlps.NewToken(tlps.RightBrace, "}", nil, 2),
				tlps.NewToken(tlps.EOF, "", nil, 1),
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
				tlps.NewToken(tlps.If, "if", nil, 1),
				tlps.NewToken(tlps.True, "true", nil, 1),
				tlps.NewToken(tlps.Colon, ":", nil, 1),
				tlps.NewToken(tlps.Newline, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 2),
				tlps.NewToken(tlps.Print, "print", nil, 2),
				tlps.NewToken(tlps.Number, "1", 1.0, 2),
				tlps.NewToken(tlps.Newline, "\n", nil, 1),
				tlps.NewToken(tlps.RightBrace, "}", nil, 2),

				// else branch
				tlps.NewToken(tlps.Else, "else", nil, 2),
				tlps.NewToken(tlps.Colon, ":", nil, 2),
				tlps.NewToken(tlps.Newline, "\n", nil, 2),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 3),
				tlps.NewToken(tlps.Print, "print", nil, 3),
				tlps.NewToken(tlps.Number, "2", 2.0, 3),
				tlps.NewToken(tlps.Newline, "\n", nil, 3),
				tlps.NewToken(tlps.RightBrace, "}", nil, 3),
				tlps.NewToken(tlps.EOF, "", nil, 3),
			},
		},
		{
			name: "for var i = 0; i < 5; i = i + 1:\n  print i\n",
			expected: []tlps.Stmt{
				tlps.NewBlock(
					[]tlps.Stmt{
						tlps.NewVar_(
							tlps.NewToken(tlps.Identifier, "i", nil, 1),
							tlps.NewLiteral(0.0),
						),
						tlps.NewWhile_(
							tlps.NewBinary(
								tlps.NewVariable(tlps.NewToken(tlps.Identifier, "i", nil, 1)),
								tlps.NewToken(tlps.Less, "<", nil, 1),
								tlps.NewLiteral(5.0),
							),
							tlps.NewBlock([]tlps.Stmt{
								tlps.NewBlock([]tlps.Stmt{
									tlps.NewPrint_(
										tlps.NewVariable(
											tlps.NewToken(tlps.Identifier, "i", nil, 2),
										),
									),
								}),
								tlps.NewExpression(
									tlps.NewAssign(
										tlps.NewToken(tlps.Identifier, "i", nil, 1),
										tlps.NewBinary(
											tlps.NewVariable(
												tlps.NewToken(tlps.Identifier, "i", nil, 1),
											),
											tlps.NewToken(tlps.Plus, "+", nil, 1),
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
				tlps.NewToken(tlps.For, "for", nil, 1),
				tlps.NewToken(tlps.Var, "var", nil, 1),
				tlps.NewToken(tlps.Identifier, "i", nil, 1),
				tlps.NewToken(tlps.Equal, "=", nil, 1),
				tlps.NewToken(tlps.Number, "0", 0.0, 1),
				tlps.NewToken(tlps.Semicolon, ";", nil, 1),
				tlps.NewToken(tlps.Identifier, "i", nil, 1),
				tlps.NewToken(tlps.Less, "<", nil, 1),
				tlps.NewToken(tlps.Number, "5", 5.0, 1),
				tlps.NewToken(tlps.Semicolon, ";", nil, 1),
				tlps.NewToken(tlps.Identifier, "i", nil, 1),
				tlps.NewToken(tlps.Equal, "=", nil, 1),
				tlps.NewToken(tlps.Identifier, "i", nil, 1),
				tlps.NewToken(tlps.Plus, "+", nil, 1),
				tlps.NewToken(tlps.Number, "1", 1.0, 1),
				tlps.NewToken(tlps.Colon, ":", nil, 1),
				tlps.NewToken(tlps.Newline, "\n", nil, 1),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 2),
				tlps.NewToken(tlps.Print, "print", nil, 2),
				tlps.NewToken(tlps.Identifier, "i", nil, 2),
				tlps.NewToken(tlps.Newline, "\n", nil, 2),
				tlps.NewToken(tlps.RightBrace, "}", nil, 2),
				tlps.NewToken(tlps.EOF, "", nil, 2),
			},
		},
	}

	// ast := tlps.NewAstPrinter()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := tlps.NewParser(r, tt.given)
			actual, _ := parser.Parse()
			// fmt.Println(ast.Print(actual))
			// fmt.Println(ast.Print(tt.expected))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
