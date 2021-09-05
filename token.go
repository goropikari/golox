package tlps

import "fmt"

// TokenType is type of type
type TokenType int

const (
	// Single-character tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Newline
	Plus
	Semicolon
	Colon
	Slash
	Star

	// One or two chacacter tokens
	Bang
	BangEqual
	BangBang
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literal
	Identifier
	String
	Number

	// keywords
	And
	Class
	Else
	Elseif
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
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
