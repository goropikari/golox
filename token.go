package tlps

import "fmt"

// TokenType is type of type
type TokenType int

const (
	// Single-character tokens
	LeftParenTT TokenType = iota
	RightParenTT
	LeftBraceTT
	RightBraceTT
	CommaTT
	DotTT
	MinusTT
	NewlineTT
	PlusTT
	SemicolonTT
	ColonTT
	SlashTT
	StarTT

	// One or two chacacter tokens
	BangTT
	BangEqualTT
	BangBangTT
	EqualTT
	EqualEqualTT
	GreaterTT
	GreaterEqualTT
	LessTT
	LessEqualTT

	// Literal
	IdentifierTT
	StringTT
	NumberTT

	// keywords
	AndTT
	ClassTT
	ElseTT
	ElseifTT
	FalseTT
	FunTT
	ForTT
	IfTT
	IncludeTT
	NilTT
	OrTT
	PassTT
	ReturnTT
	SuperTT
	ThisTT
	TrueTT
	VarTT
	WhileTT

	EOFTT
)

// Token is struct of token
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// TokenList is slice of Token
type TokenList []*Token

// NewToken is constructor of Token
func NewToken(tt TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tt,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

// String stringfy Token
func (t *Token) String() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v", t.Type, t.Lexeme, t.Literal, t.Line)
}
