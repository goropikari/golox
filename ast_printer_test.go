package mylang_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/mylang"
)

func TestAstPrinter(t *testing.T) {
	expression := mylang.NewBinary(
		mylang.NewUnary(mylang.NewToken(mylang.Minus, "-", nil, 1), mylang.NewLiteral(123)),
		mylang.NewToken(mylang.Star, "*", nil, 1),
		mylang.NewGrouping(mylang.NewLiteral(45.67)),
	)
	fmt.Println(mylang.NewAstPrinter().Print(expression))
}
