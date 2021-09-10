package tlps

import (
	"bytes"
	"strconv"
	"unicode"
)

// Scanner is struct of Scanner
type Scanner struct {
	keywords    map[string]TokenType
	indent      *IndentStack
	isFirst     bool // use when count indentation level
	runtime     *Runtime
	source      *bytes.Buffer
	sourceRunes []rune
	tokens      TokenList
	start       int
	current     int
	line        int
}

// NewScanner is constructor of Scanner
func NewScanner(r *Runtime, b *bytes.Buffer) *Scanner {
	var keywords = map[string]TokenType{
		"and":    AndTT,
		"class":  ClassTT,
		"else":   ElseTT,
		"elseif": ElseifTT,
		"false":  FalseTT,
		"for":    ForTT,
		"fun":    FunTT,
		"if":     IfTT,
		"nil":    NilTT,
		"or":     OrTT,
		"print":  PrintTT,
		"return": ReturnTT,
		"super":  SuperTT,
		"this":   ThisTT,
		"true":   TrueTT,
		"var":    VarTT,
		"while":  WhileTT,
	}

	indent := NewIndentStack()
	indent.Push(0)

	return &Scanner{
		keywords:    keywords,
		indent:      indent,
		isFirst:     true,
		runtime:     r,
		source:      b,
		sourceRunes: bytes.Runes(b.Bytes()),
		tokens:      []*Token{},
		start:       0,
		current:     0,
		line:        1,
	}
}

// ScanTokens generates tokens from given source code.
func (s *Scanner) ScanTokens() TokenList {
	for !s.isAtEnd() {
		s.addBlock()
		s.start = s.current
		s.scanToken()
		s.removeUselessNewline()
	}

	for s.indent.Peek() != 0 {
		s.indent.Pop()
		s.tokens = append(s.tokens, NewToken(RightBraceTT, "}", nil, s.line))
	}

	s.tokens = append(s.tokens, NewToken(EOFTT, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c, _, _ := s.advance()

	switch c {
	case '(':
		s.addToken(LeftParenTT, nil)
		break
	case ')':
		s.addToken(RightParenTT, nil)
		break
	// case '{':
	// 	s.addToken(LeftBraceTT, nil)
	// 	break
	// case '}':
	// 	s.addToken(RightBraceTT, nil)
	// 	break
	case ',':
		s.addToken(CommaTT, nil)
		break
	case '.':
		s.addToken(DotTT, nil)
		break
	case '-':
		s.addToken(MinusTT, nil)
		break
	case '+':
		s.addToken(PlusTT, nil)
		break
	case ';':
		s.addToken(SemicolonTT, nil)
		break
	case ':':
		s.addToken(ColonTT, nil)
		break
	case '*':
		s.addToken(StarTT, nil)
		break
	case '!':
		var tt TokenType
		if s.match('=') {
			tt = BangEqualTT
		} else {
			tt = BangTT
		}
		s.addToken(tt, nil)
	case '=':
		var tt TokenType
		if s.match('=') {
			tt = EqualEqualTT
		} else {
			tt = EqualTT
		}
		s.addToken(tt, nil)
	case '<':
		var tt TokenType
		if s.match('=') {
			tt = LessEqualTT
		} else {
			tt = LessTT
		}
		s.addToken(tt, nil)
	case '>':
		var tt TokenType
		if s.match('=') {
			tt = GreaterEqualTT
		} else {
			tt = GreaterTT
		}
		s.addToken(tt, nil)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SlashTT, nil)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.addNewline()
		break
	case '"':
		s.addString()
		break
	default:
		if unicode.IsDigit(c) {
			s.addNumber()
		} else if unicode.IsLetter(c) {
			s.addIdentifier()
		} else {
			s.runtime.ErrorMessage(s.line, "Unexpected character.")
		}
		break
	}
}

func (s *Scanner) removeUselessNewline() {
	tokens := make([]*Token, 0)
	isPreviousNewline := true
	for _, token := range s.tokens {
		if token.Type == NewlineTT {
			if !isPreviousNewline {
				tokens = append(tokens, token)
			}
			isPreviousNewline = true
		} else {
			isPreviousNewline = false
			tokens = append(tokens, token)
		}
	}

	s.tokens = tokens
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.peek() != expected {
		return false
	}

	s.advance()
	return true

}

func (s *Scanner) advance() (rune, int, error) {
	r, size, err := s.source.ReadRune()
	s.current++

	return r, size, err
}

func (s *Scanner) addToken(tt TokenType, literal interface{}) {
	s.isFirst = false
	text := string(s.sourceRunes[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(tt, text, literal, s.line))
}

func (s *Scanner) addNewline() {
	s.tokens = append(s.tokens, NewToken(NewlineTT, "\\n", nil, s.line))
	s.line++
	s.isFirst = true
}

func (s *Scanner) addBlock() {
	if !s.isFirst {
		return
	}

	depth := 0
	for s.match(' ') {
		depth++
	}

	// skip comment or empty line
	if s.peek() == '\n' || (s.peek() == '/' && s.peekNext() == '/') {
		return
	}

	d := s.indent.Peek()
	if d < depth {
		s.indent.Push(depth)
		s.tokens = append(s.tokens, NewToken(LeftBraceTT, "{", nil, s.line))
	} else if d > depth {
		cnt := 0
		for s.indent.Pop() != -1 {
			cnt++
			if s.indent.Peek() == depth {
				break
			}
		}

		if s.indent.IsEmpty() {
			s.runtime.ErrorMessage(s.line, "unindent does not match any outer indentation level")
		}

		for i := 0; i < cnt; i++ {
			s.tokens = append(s.tokens, NewToken(RightBraceTT, "}", nil, s.line))
		}
	}
}

func (s *Scanner) addString() {
	isEscape := false // define isEscape to handle \"
	for (isEscape || s.peek() != '"') && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		if s.peek() == '\\' {
			isEscape = !isEscape
		} else {
			isEscape = false
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.runtime.ErrorMessage(s.line, "Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	// value := string(s.sourceRunes[s.start+1 : s.current-1])
	value := s.handleEscapeCharacter()
	s.addToken(StringTT, string(value))
}

func (s *Scanner) addNumber() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	f, _ := strconv.ParseFloat(string(s.sourceRunes[s.start:s.current]), 64)
	s.addToken(NumberTT, f)
}

func (s *Scanner) addIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.sourceRunes[s.start:s.current])
	typ, ok := s.keywords[text]
	if ok {
		s.addToken(typ, nil)
		return
	}

	s.addToken(IdentifierTT, nil)
}

func isAlphaNumeric(x rune) bool {
	return unicode.IsDigit(x) || unicode.IsLetter(x) || x == '_'
}

func (s *Scanner) isAtEnd() bool {
	return s.source.Len() == 0
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.sourceRunes[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.sourceRunes) {
		return 0
	}
	return s.sourceRunes[s.current+1]
}

func (s *Scanner) handleEscapeCharacter() string {
	runes := s.sourceRunes[s.start+1 : s.current-1]
	value := make([]rune, 0)
	var prevc rune
	for _, v := range runes {
		var c rune
		if prevc == '\\' {
			// https://en.wikipedia.org/wiki/C_syntax#Backslash_escapes
			switch v {
			case '\\':
				c = v
			case '"':
				c = v
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 'b':
				c = '\b'
			case 't':
				c = '\t'
			case 'f':
				c = '\f'
			case 'v':
				c = '\v'
			default:
				s.runtime.ErrorMessage(s.line, "invalid escape sequence")
				return ""
			}

			if v == '\\' {
				prevc = 0
			} else {
				prevc = v
			}
			value = append(value, c)
		} else if v == '\\' {
			prevc = v
		} else {
			prevc = v
			value = append(value, v)
		}
	}

	return string(value)
}
