package golox_test

import (
	"testing"

	"github.com/goropikari/golox"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	r := golox.NewRuntime()

	var tests = []struct {
		name     string
		expected interface{}
		given    []golox.Stmt
	}{
		{
			name:     "1.3 + 1.2",
			expected: "2.5",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(golox.NewLiteral(1.3), golox.NewToken(golox.PlusTT, "+", nil, 1), golox.NewLiteral(1.2))),
			},
		},
		{
			name:     "1.3 * 1.2",
			expected: "1.56",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(golox.NewLiteral(1.3), golox.NewToken(golox.StarTT, "*", nil, 1), golox.NewLiteral(1.2))),
			},
		},
		{
			name:     "2 / 4",
			expected: "0.5",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(golox.NewLiteral(2.0), golox.NewToken(golox.SlashTT, "/", nil, 1), golox.NewLiteral(4.0))),
			},
		},
		{
			name:     "string + string",
			expected: "foo bar",
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(golox.NewLiteral("foo "), golox.NewToken(golox.PlusTT, "+", nil, 1), golox.NewLiteral("bar"))),
			},
		},
		{
			name:     "function",
			expected: "13",
			given: []golox.Stmt{
				// fun f(x, y):
				//   return x + y
				// f(11, 2)
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
								golox.NewToken(golox.PlusTT, "+", nil, 1),
								golox.NewVariable(golox.NewToken(golox.IdentifierTT, "y", nil, 2)),
							),
						),
					},
				),
				golox.NewExpression(
					golox.NewCall(
						golox.NewVariable(golox.NewToken(golox.IdentifierTT, "f", nil, 3)),
						golox.NewToken(golox.LeftParenTT, "(", nil, 3),
						[]golox.Expr{
							golox.NewLiteral(float64(11)),
							golox.NewLiteral(float64(2)),
						},
					),
				),
			},
		},
		{
			name:     "class",
			expected: "hoge",
			given: []golox.Stmt{
				// class Hoge:
				//   init(x):
				//     this.x = x
				// Hoge("hoge").x
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
				golox.NewExpression(
					golox.NewGet(
						golox.NewCall(
							golox.NewVariable(
								golox.NewToken(golox.IdentifierTT, "Hoge", nil, 4),
							),
							golox.NewToken(golox.LeftParenTT, "(", nil, 4),
							[]golox.Expr{
								golox.NewLiteral("hoge"),
							},
						),
						golox.NewToken(golox.IdentifierTT, "x", nil, 4),
					),
				),
			},
		},
		// {
		// 	name:     "super class",
		// 	expected: "hoge",
		// 	given: []golox.Stmt{
		// 		// class Hoge:
		// 		//   init(x):
		// 		//     this.x = x
		// 		// class Piyo(Hoge):
		// 		//   pass
		// 		// Piyo("hoge").x
		// 		golox.NewClass(
		// 			golox.NewToken(golox.IdentifierTT, "Hoge", nil, 1),
		// 			nil,
		// 			[]*golox.Function{
		// 				golox.NewFunction(
		// 					golox.NewToken(golox.IdentifierTT, "init", nil, 2),
		// 					[]*golox.Token{
		// 						golox.NewToken(golox.IdentifierTT, "x", nil, 2),
		// 					},
		// 					[]golox.Stmt{
		// 						golox.NewExpression(
		// 							golox.NewSet(
		// 								golox.NewThis(
		// 									golox.NewToken(golox.ThisTT, "this", nil, 3),
		// 								),
		// 								golox.NewToken(golox.IdentifierTT, "x", nil, 3),
		// 								golox.NewVariable(
		// 									golox.NewToken(golox.IdentifierTT, "x", nil, 3),
		// 								),
		// 							),
		// 						),
		// 					},
		// 				).(*golox.Function),
		// 			},
		// 		),
		// 		golox.NewClass(
		// 			golox.NewToken(golox.IdentifierTT, "Piyo", nil, 4),
		// 			golox.NewVariable(
		// 				golox.NewToken(golox.IdentifierTT, "Hoge", nil, 4),
		// 			).(*golox.Variable),
		// 			[]*golox.Function{
		// 				golox.NewFunction(
		// 					golox.NewToken(golox.PassTT, "pass", nil, 5),
		// 					[]*golox.Token{},
		// 					[]golox.Stmt{},
		// 				).(*golox.Function),
		// 			},
		// 		),
		// 		golox.NewExpression(
		// 			golox.NewGet(
		// 				golox.NewCall(
		// 					golox.NewVariable(
		// 						golox.NewToken(golox.IdentifierTT, "Piyo", nil, 6),
		// 					),
		// 					golox.NewToken(golox.LeftParenTT, "(", nil, 6),
		// 					[]golox.Expr{
		// 						golox.NewLiteral("hoge"),
		// 					},
		// 				),
		// 				golox.NewToken(golox.IdentifierTT, "x", nil, 6),
		// 			),
		// 		),
		// 	},
		// },
		{
			name:     "scope resolution",
			expected: "20",
			given: []golox.Stmt{
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
				golox.NewFunction(
					golox.NewToken(golox.IdentifierTT, "f", nil, 1),
					[]*golox.Token{},
					[]golox.Stmt{
						golox.NewVar(
							golox.NewToken(golox.IdentifierTT, "a", nil, 2),
							golox.NewLiteral(10.0),
						),

						golox.NewFunction(
							golox.NewToken(golox.IdentifierTT, "g", nil, 3),
							[]*golox.Token{},
							[]golox.Stmt{
								golox.NewFunction(
									golox.NewToken(golox.IdentifierTT, "h", nil, 4),
									[]*golox.Token{},
									[]golox.Stmt{
										golox.NewReturn(
											golox.NewToken(golox.ReturnTT, "return", nil, 2),
											golox.NewVariable(golox.NewToken(golox.IdentifierTT, "a", nil, 2)),
										),
									},
								),
								golox.NewVar(
									golox.NewToken(golox.IdentifierTT, "x", nil, 6),
									golox.NewCall(
										golox.NewVariable(golox.NewToken(golox.IdentifierTT, "h", nil, 6)),
										golox.NewToken(golox.LeftParenTT, "(", nil, 6),
										[]golox.Expr{},
									),
								),
								golox.NewVar(
									golox.NewToken(golox.IdentifierTT, "a", nil, 7),
									golox.NewLiteral(123.0),
								),
								golox.NewVar(
									golox.NewToken(golox.IdentifierTT, "y", nil, 8),
									golox.NewCall(
										golox.NewVariable(golox.NewToken(golox.IdentifierTT, "h", nil, 8)),
										golox.NewToken(golox.LeftParenTT, "(", nil, 8),
										[]golox.Expr{},
									),
								),
								golox.NewReturn(
									golox.NewToken(golox.ReturnTT, "return", nil, 9),
									golox.NewBinary(
										golox.NewVariable(golox.NewToken(golox.IdentifierTT, "x", nil, 9)),
										golox.NewToken(golox.PlusTT, "+", nil, 9),
										golox.NewVariable(golox.NewToken(golox.IdentifierTT, "y", nil, 9)),
									),
								),
							},
						),
						golox.NewReturn(
							golox.NewToken(golox.ReturnTT, "return", nil, 10),
							golox.NewVariable(golox.NewToken(golox.IdentifierTT, "g", nil, 10)),
						),
					},
				),
				golox.NewVar(
					golox.NewToken(golox.IdentifierTT, "fn", nil, 11),
					golox.NewCall(
						golox.NewVariable(golox.NewToken(golox.IdentifierTT, "f", nil, 11)),
						golox.NewToken(golox.LeftParenTT, "(", nil, 11),
						[]golox.Expr{},
					),
				),
				golox.NewExpression(
					golox.NewCall(
						golox.NewVariable(golox.NewToken(golox.IdentifierTT, "fn", nil, 12)),
						golox.NewToken(golox.LeftParenTT, "(", nil, 12),
						[]golox.Expr{},
					),
				),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interpreter := golox.NewInterpreter(r)
			resolver := golox.NewResolver(r, interpreter)
			resolver.ResolveStmts(tt.given)
			actual, _ := interpreter.Interpret(tt.given)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInterpreter_Error(t *testing.T) {
	r := golox.NewRuntime()
	plus := golox.NewToken(golox.PlusTT, "+", nil, 1)

	var tests = []struct {
		name     string
		expected interface{}
		err      error
		given    []golox.Stmt
	}{
		{
			name:     "number + string",
			expected: "nil",
			err:      golox.RuntimeError.New(plus, "Operands must be two numbers or two strings."),
			given: []golox.Stmt{
				golox.NewExpression(golox.NewBinary(golox.NewLiteral(1.5), plus, golox.NewLiteral("bar"))),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interpreter := golox.NewInterpreter(r)
			actual, err := interpreter.Interpret(tt.given)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.err, err)
		})
	}
}
