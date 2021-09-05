package mylang_test

import (
	"testing"

	"github.com/goropikari/mylang"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	r := mylang.NewRuntime()

	var tests = []struct {
		name     string
		expected string
		given    mylang.TokenList
	}{
		{
			name:     "+ * precedence: 1 + 2 * 3",
			expected: "(+ 1 (* 2 3))",
			given: mylang.TokenList{
				mylang.NewToken(mylang.Number, "1", 1.0, 1),
				mylang.NewToken(mylang.Plus, "+", nil, 1),
				mylang.NewToken(mylang.Number, "2", 2.0, 1),
				mylang.NewToken(mylang.Star, "*", nil, 1),
				mylang.NewToken(mylang.Number, "3", 3.0, 1),
				mylang.NewToken(mylang.EOF, "", nil, 1),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := mylang.NewParser(r, tt.given)
			expression := parser.Parse()
			actual := mylang.NewAstPrinter().Print(expression)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
