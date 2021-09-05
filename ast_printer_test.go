package tlps_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/tlps"
)

func TestAstPrinter(t *testing.T) {
	expression := tlps.NewBinary(
		tlps.NewUnary(tlps.NewToken(tlps.Minus, "-", nil, 1), tlps.NewLiteral(123)),
		tlps.NewToken(tlps.Star, "*", nil, 1),
		tlps.NewGrouping(tlps.NewLiteral(45.67)),
	)
	fmt.Println(tlps.NewAstPrinter().Print(expression))
}
