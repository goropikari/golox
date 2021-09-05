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
			name:     "",
			expected: "(* (- 123) (group 45.67))",
			given: []tlps.Stmt{
				tlps.NewExpression(tlps.NewBinary(
					tlps.NewUnary(tlps.NewToken(tlps.Minus, "-", nil, 1), tlps.NewLiteral(123)),
					tlps.NewToken(tlps.Star, "*", nil, 1),
					tlps.NewGrouping(tlps.NewLiteral(45.67))),
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
