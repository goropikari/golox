package tlps_test

import (
	"testing"

	"github.com/goropikari/tlps"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	r := tlps.NewRuntime()

	var tests = []struct {
		name     string
		expected interface{}
		given    []tlps.Stmt
	}{
		{
			name:     "1.3 + 1.2",
			expected: "2.5",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(1.3), tlps.NewToken(tlps.PlusTT, "+", nil, 1), tlps.NewLiteral(1.2))),
			},
		},
		{
			name:     "1.3 * 1.2",
			expected: "1.56",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(1.3), tlps.NewToken(tlps.StarTT, "*", nil, 1), tlps.NewLiteral(1.2))),
			},
		},
		{
			name:     "2 / 4",
			expected: "0.5",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(2.0), tlps.NewToken(tlps.SlashTT, "/", nil, 1), tlps.NewLiteral(4.0))),
			},
		},
		{
			name:     "string + string",
			expected: "foo bar",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral("foo "), tlps.NewToken(tlps.PlusTT, "+", nil, 1), tlps.NewLiteral("bar"))),
			},
		},
		{
			name:     "function",
			expected: "13",
			given: []tlps.Stmt{
				// fun f(x, y):
				//   return x + y
				// f(11, 2)
				tlps.NewFunction(
					tlps.NewToken(tlps.IdentifierTT, "f", nil, 1),
					[]*tlps.Token{
						tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
						tlps.NewToken(tlps.IdentifierTT, "y", nil, 1),
					},
					[]tlps.Stmt{
						tlps.NewReturn(
							tlps.NewToken(tlps.ReturnTT, "return", nil, 2),
							tlps.NewBinary(
								tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "x", nil, 2)),
								tlps.NewToken(tlps.PlusTT, "+", nil, 1),
								tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "y", nil, 2)),
							),
						),
					},
				),
				tlps.NewExpression(
					tlps.NewCall(
						tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "f", nil, 3)),
						tlps.NewToken(tlps.LeftParenTT, "(", nil, 3),
						[]tlps.Expr{
							tlps.NewLiteral(float64(11)),
							tlps.NewLiteral(float64(2)),
						},
					),
				),
			},
		},

		{
			name:     "class",
			expected: "hoge",
			given: []tlps.Stmt{
				// class Hoge:
				//   init(x):
				//     this.x = x
				// Hoge("hoge").x
				tlps.NewClass(
					tlps.NewToken(tlps.IdentifierTT, "Hoge", nil, 1),
					[]*tlps.Function{
						tlps.NewFunction(
							tlps.NewToken(tlps.IdentifierTT, "init", nil, 2),
							[]*tlps.Token{
								tlps.NewToken(tlps.IdentifierTT, "x", nil, 2),
							},
							[]tlps.Stmt{
								tlps.NewExpression(
									tlps.NewSet(
										tlps.NewThis(
											tlps.NewToken(tlps.ThisTT, "this", nil, 3),
										),
										tlps.NewToken(tlps.IdentifierTT, "x", nil, 3),
										tlps.NewVariable(
											tlps.NewToken(tlps.IdentifierTT, "x", nil, 3),
										),
									),
								),
							},
						).(*tlps.Function),
					},
				),
				tlps.NewExpression(
					tlps.NewGet(
						tlps.NewCall(
							tlps.NewVariable(
								tlps.NewToken(tlps.IdentifierTT, "Hoge", nil, 4),
							),
							tlps.NewToken(tlps.LeftParenTT, "(", nil, 4),
							[]tlps.Expr{
								tlps.NewLiteral("hoge"),
							},
						),
						tlps.NewToken(tlps.IdentifierTT, "x", nil, 4),
					),
				),
			},
		},
		{
			name:     "scope resolution",
			expected: "20",
			given: []tlps.Stmt{
				// fun f():
				//   var a = 10
				//   fun g():
				//     fun h():
				//       return a
				//     var x = h()
				//     var a = 123
				//     var y = h()
				//     return x + y
				//   return g
				// var fn = f()
				// fn()
				tlps.NewFunction(
					tlps.NewToken(tlps.IdentifierTT, "f", nil, 1),
					[]*tlps.Token{},
					[]tlps.Stmt{
						tlps.NewVar(
							tlps.NewToken(tlps.IdentifierTT, "a", nil, 2),
							tlps.NewLiteral(10.0),
						),

						tlps.NewFunction(
							tlps.NewToken(tlps.IdentifierTT, "g", nil, 3),
							[]*tlps.Token{},
							[]tlps.Stmt{
								tlps.NewFunction(
									tlps.NewToken(tlps.IdentifierTT, "h", nil, 4),
									[]*tlps.Token{},
									[]tlps.Stmt{
										tlps.NewReturn(
											tlps.NewToken(tlps.ReturnTT, "return", nil, 2),
											tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "a", nil, 2)),
										),
									},
								),
								tlps.NewVar(
									tlps.NewToken(tlps.IdentifierTT, "x", nil, 6),
									tlps.NewCall(
										tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "h", nil, 6)),
										tlps.NewToken(tlps.LeftParenTT, "(", nil, 6),
										[]tlps.Expr{},
									),
								),
								tlps.NewVar(
									tlps.NewToken(tlps.IdentifierTT, "a", nil, 7),
									tlps.NewLiteral(123.0),
								),
								tlps.NewVar(
									tlps.NewToken(tlps.IdentifierTT, "y", nil, 8),
									tlps.NewCall(
										tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "h", nil, 8)),
										tlps.NewToken(tlps.LeftParenTT, "(", nil, 8),
										[]tlps.Expr{},
									),
								),
								tlps.NewReturn(
									tlps.NewToken(tlps.ReturnTT, "return", nil, 9),
									tlps.NewBinary(
										tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "x", nil, 9)),
										tlps.NewToken(tlps.PlusTT, "+", nil, 9),
										tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "y", nil, 9)),
									),
								),
							},
						),
						tlps.NewReturn(
							tlps.NewToken(tlps.ReturnTT, "return", nil, 10),
							tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "g", nil, 10)),
						),
					},
				),
				tlps.NewVar(
					tlps.NewToken(tlps.IdentifierTT, "fn", nil, 11),
					tlps.NewCall(
						tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "f", nil, 11)),
						tlps.NewToken(tlps.LeftParenTT, "(", nil, 11),
						[]tlps.Expr{},
					),
				),
				tlps.NewExpression(
					tlps.NewCall(
						tlps.NewVariable(tlps.NewToken(tlps.IdentifierTT, "fn", nil, 12)),
						tlps.NewToken(tlps.LeftParenTT, "(", nil, 12),
						[]tlps.Expr{},
					),
				),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interpreter := tlps.NewInterpreter(r)
			resolver := tlps.NewResolver(r, interpreter)
			resolver.ResolveStmts(tt.given)
			actual, _ := interpreter.Interpret(tt.given)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInterpreter_Error(t *testing.T) {
	r := tlps.NewRuntime()
	plus := tlps.NewToken(tlps.PlusTT, "+", nil, 1)

	var tests = []struct {
		name     string
		expected interface{}
		err      error
		given    []tlps.Stmt
	}{
		{
			name:     "number + string",
			expected: "nil",
			err:      tlps.RuntimeError.New(plus, "Operands must be two numbers or two strings."),
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(1.5), plus, tlps.NewLiteral("bar"))),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interpreter := tlps.NewInterpreter(r)
			actual, err := interpreter.Interpret(tt.given)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.err, err)
		})
	}
}
