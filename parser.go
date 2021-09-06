package tlps

// Parser is struct of parser
type Parser struct {
	runtime *Runtime
	tokens  TokenList
	current int
}

// NewParser is constructor of Parser
func NewParser(runtime *Runtime, tokens TokenList) *Parser {
	return &Parser{
		runtime: runtime,
		tokens:  tokens,
		current: 0,
	}
}

// Parse parses given tokens
func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(Var) {
		stmt, err := p.varDecralation()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		return stmt, nil
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
	}
	return stmt, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(LeftBrace) {
		b, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(b), nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	var initializer Stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		init, err := p.varDecralation()
		if err != nil {
			return nil, err
		}
		initializer = init
	} else {
		init, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = init
	}

	var condition Expr
	if !p.check(Semicolon) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}
		condition = cond
	}
	p.consume(Semicolon, "Expect ';' after loop condition")

	var increment Expr
	if !p.check(Colon) {
		inc, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = inc
	}
	p.consume(Colon, "Expect ':' after for clauses.")
	p.consume(Newline, "Expect '\\n' after for clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)})
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile_(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(Colon, "Expect ':' after if condition.")
	p.consume(Newline, "Expect '\\n' after if condition")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if p.match(Elseif) {
		elseBranch, err = p.ifStatement()
		if err != nil {
			return nil, err
		}
	}
	if p.match(Else) {
		p.consume(Colon, "Expect ':' after else.")
		p.consume(Newline, "Expect '\\n' after else")
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIf_(condition, thenBranch, elseBranch), nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	// _, err = p.consume(Semicolon, "Expect ';' after value.")
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewPrint_(value), nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(Colon, "Expect ':' after condition")
	p.consume(Newline, "Expect '\\n' after condition")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return NewWhile_(condition, body), nil
}

func (p *Parser) varDecralation() (Stmt, error) {
	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	// _, err = p.consume(Semicolon, "Expect ';' after variable declaration")
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewVar_(name, initializer), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	// _, err = p.consume(Semicolon, "Expect ';' after expression")
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewExpression(expr), nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	p.consume(RightBrace, "Expect '}' after block.")
	return statements, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if expr.IsType(&Variable{}) {
			name := expr.(*Variable).Name
			return NewAssign(name, value), nil
		}

		p.runtime.ErrorTokenMessage(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return NewLiteral(false), nil
	}
	if p.match(True) {
		return NewLiteral(true), nil
	}
	if p.match(Nil) {
		return NewLiteral(nil), nil
	}

	if p.match(Number, String) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.match(Identifier) {
		return NewVariable(p.previous()), nil
	}

	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return NewGrouping(expr), nil
	}

	if p.match(Newline) {
		return NewLiteral('\n'), nil
	}

	return nil, p.NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(typ TokenType, message string) (*Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}

	return nil, p.NewParseError(p.peek(), message)
}

func (p *Parser) consumeTerm() (*Token, error) {
	if p.check(Newline) {
		return p.advance(), nil
	}
	if p.check(Semicolon) {
		t := p.advance()
		if p.check(Newline) {
			p.advance()
		}
		return t, nil
	}

	return nil, p.NewParseError(p.peek(), "Expect ';' or '\\n' after expression")
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		switch p.previous().Type {
		case Semicolon:
		case Newline:
			return
		}
		// if p.previous().Type == Semicolon {
		// 	return
		// }

		switch p.peek().Type {
		case Class:
		case Fun:
		case Var:
		case For:
		case If:
		case While:
		case Print:
		case Return:
			return
		}

		p.advance()
	}
}

// NewParseError is constructor of ParseError
func (p *Parser) NewParseError(token *Token, message string) error {
	p.runtime.ErrorTokenMessage(token, message)
	return ParseError.New(token, message)
}
