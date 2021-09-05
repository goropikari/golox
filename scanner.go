package tlps

import (
	"bytes"
	"strconv"
	"unicode"

	"github.com/goropikari/tlps/collections/stack"
)

// Scanner is struct of Scanner
type Scanner struct {
	keywords    map[string]TokenType
	stack_      *stack.Stack
	isFirst     bool
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
	st := stack.NewStack()
	st.Push(0)

	var keywords = map[string]TokenType{
		"and":    And,
		"class":  Class,
		"else":   Else,
		"elseif": Elseif,
		"false":  False,
		"for":    For,
		"fun":    Fun,
		"if":     If,
		"nil":    Nil,
		"or":     Or,
		"print":  Print,
		"return": Return,
		"super":  Super,
		"this":   This,
		"true":   True,
		"var":    Var,
		"while":  While,
	}

	return &Scanner{
		keywords:    keywords,
		stack_:      st,
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
	}

	for s.stack_.Top() != 0 {
		s.stack_.Pop()
		s.tokens = append(s.tokens, NewToken(RightBrace, "}", nil, s.line))
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c, _, _ := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen, nil)
		break
	case ')':
		s.addToken(RightParen, nil)
		break
	// case '{':
	// 	s.addToken(LeftBrace, nil)
	// 	break
	// case '}':
	// 	s.addToken(RightBrace, nil)
	// 	break
	case ',':
		s.addToken(Comma, nil)
		break
	case '.':
		s.addToken(Dot, nil)
		break
	case '-':
		s.addToken(Minus, nil)
		break
	case '+':
		s.addToken(Plus, nil)
		break
	case ';':
		s.addToken(Semicolon, nil)
		break
	case ':':
		s.addToken(Colon, nil)
		break
	case '*':
		s.addToken(Star, nil)
		break
	case '!':
		var tt TokenType
		if s.match('=') {
			tt = BangEqual
		} else {
			tt = Bang
		}
		s.addToken(tt, nil)
	case '=':
		var tt TokenType
		if s.match('=') {
			tt = EqualEqual
		} else {
			tt = Equal
		}
		s.addToken(tt, nil)
	case '<':
		var tt TokenType
		if s.match('=') {
			tt = LessEqual
		} else {
			tt = Less
		}
		s.addToken(tt, nil)
	case '>':
		var tt TokenType
		if s.match('=') {
			tt = GreaterEqual
		} else {
			tt = Greater
		}
		s.addToken(tt, nil)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash, nil)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
		s.isFirst = true
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

	d := s.stack_.Top()
	if d < depth {
		s.stack_.Push(depth)
		s.tokens = append(s.tokens, NewToken(LeftBrace, "{", nil, s.line))
	} else if d > depth {
		cnt := 0
		for s.stack_.Pop() != -1 {
			cnt++
			if s.stack_.Top() == depth {
				break
			}
		}

		if s.stack_.IsEmpty() {
			s.runtime.ErrorMessage(s.line, "unindent does not match any outer indentation level")
		}

		for i := 0; i < cnt; i++ {
			s.tokens = append(s.tokens, NewToken(RightBrace, "}", nil, s.line))
		}
	}
}

func (s *Scanner) addString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.match('\n') {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.runtime.ErrorMessage(s.line, "Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	value := s.sourceRunes[s.start+1 : s.current-1]
	s.addToken(String, value)
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
	s.addToken(Number, f)
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

	s.addToken(Identifier, nil)
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
