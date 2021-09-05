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
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(1.3), tlps.NewToken(tlps.Plus, "+", nil, 1), tlps.NewLiteral(1.2))),
			},
		},
		{
			name:     "1.3 * 1.2",
			expected: "1.56",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(1.3), tlps.NewToken(tlps.Star, "*", nil, 1), tlps.NewLiteral(1.2))),
			},
		},
		{
			name:     "2 / 4",
			expected: "0.5",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral(2.0), tlps.NewToken(tlps.Slash, "/", nil, 1), tlps.NewLiteral(4.0))),
			},
		},
		{
			name:     "string + string",
			expected: "foo bar",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(tlps.NewLiteral("foo "), tlps.NewToken(tlps.Plus, "+", nil, 1), tlps.NewLiteral("bar"))),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			interpreter := tlps.NewInterpreter(r)
			actual, _ := interpreter.Interpret(tt.given)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInterpreter_Error(t *testing.T) {
	r := tlps.NewRuntime()
	plus := tlps.NewToken(tlps.Plus, "+", nil, 1)

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
