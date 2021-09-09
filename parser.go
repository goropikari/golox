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
	if p.match(Fun) {
		return p.function("function")
	}
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
	if p.match(Return) {
		return p.returnStatement()
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
	_, err := p.consume(Semicolon, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(Colon) {
		inc, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = inc
	}
	_, err = p.consume(Colon, "Expect ':' after for clauses.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Newline, "Expect '\\n' after for clauses.")
	if err != nil {
		return nil, err
	}

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
	_, err = p.consume(Colon, "Expect ':' after if condition.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Newline, "Expect '\\n' after if condition")
	if err != nil {
		return nil, err
	}

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
		_, err := p.consume(Colon, "Expect ':' after else.")
		if err != nil {
			return nil, err
		}
		_, err = p.consume(Newline, "Expect '\\n' after else")
		if err != nil {
			return nil, err
		}
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
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewPrint_(value), nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr = nil
	if !p.check(Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consumeTerm()
	if err != nil {
		return nil, err
	}
	return NewReturn_(keyword, value), nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Colon, "Expect ':' after condition")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Newline, "Expect '\\n' after condition")
	if err != nil {
		return nil, err
	}
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

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(Identifier, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftParen, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}
	parameters := make([]*Token, 0)

	if !p.check(RightParen) {
		for {
			if len(parameters) >= 255 {
				return nil, p.NewParseError(p.peek(), "Can't have more than 255 parameters.")
			}

			token, err := p.consume(Identifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, token)

			if !p.match(Comma) {
				break
			}
		}
	}

	_, err = p.consume(RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Colon, "Expect ':' after '('")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Newline, "Expect '\\n' before "+kind+" body.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftBrace, "Expected an indented block as "+kind+" body.")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return NewFunction(name, parameters, body), nil
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

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)
	if !p.check(RightParen) {
		for {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			if len(arguments) >= 255 {
				return nil, p.NewParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, expr)

			if !p.match(Comma) {
				break
			}
		}
	}

	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return NewCall(callee, paren, arguments), nil
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
			return
		case Newline:
			return
		}
		// if p.previous().Type == Semicolon {
		// 	return
		// }

		switch p.peek().Type {
		case Class:
			return
		case Fun:
			return
		case Var:
			return
		case For:
			return
		case If:
			return
		case While:
			return
		case Print:
			return
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
