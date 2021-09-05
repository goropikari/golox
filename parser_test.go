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
		expected string
		given    tlps.TokenList
	}{
		{
			name:     "+ * precedence: 1 + 2 * 3",
			expected: "(+ 1 (* 2 3))",
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := tlps.NewParser(r, tt.given)
			expression, _ := parser.Parse()
			actual, _ := tlps.NewAstPrinter().Print(expression)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
