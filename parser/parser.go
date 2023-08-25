package parser

import (
	"github.com/arturoeanton/go-r2lox/lexer"
	"github.com/arturoeanton/go-r2lox/r2loxerrors"
)

type Parser struct {
	Tokens  []lexer.Token
	Current int
	Start   int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
		Start:   0,
	}
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) ReturnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(lexer.SEMICOLON) {
		value = p.Expression()
	}

	p.consume(lexer.SEMICOLON, "Expect ';' after return value.")
	return Return{Keyword: keyword, Value: value}
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right := p.Unary()
		return Unary{Operator: operator, Right: right}
	}

	return p.Call()
}

func (p *Parser) errors(token lexer.Token, message string) {
	if token.Type == lexer.EOF {
		r2loxerrors.Errors(token.Line, " at end "+message)
	} else {
		r2loxerrors.Errors(token.Line, " at '"+token.Lexeme+"' "+message)
	}
}

func (p *Parser) Primary() Expr {
	if p.match(lexer.FALSE) {
		return Literal{Value: false}
	}
	if p.match(lexer.TRUE) {
		return Literal{Value: true}
	}
	if p.match(lexer.NIL) {
		return Literal{Value: nil}
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
		return Literal{Value: p.previous().Literal}
	}

	if p.match(lexer.IDENTIFIER) {
		return Var{Name: p.previous()}
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.Equality()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{Expression: expr}
	}

	if p.match(lexer.EOF) {
		return Literal{Value: nil}
	}

	return nil
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) peek() lexer.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() lexer.Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) consume(t lexer.TokenType, message string) lexer.Token {
	if p.check(t) {
		return p.advance()
	}

	panic(message)
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr, ok := expr.(Var); ok {
			name := expr.Name
			return Assign{Name: name, Value: value}
		}

		r2loxerrors.Errors(equals.Line, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(lexer.OR) {
		operator := p.previous()
		right := p.and()
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.Equality()

	for p.match(lexer.AND) {
		operator := p.previous()
		right := p.Equality()
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == lexer.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) Declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()

	if p.match(lexer.FUN) {
		return p.Function("function")
	}

	if p.match(lexer.VAR) {
		return p.VarDeclaration()
	}

	return p.Statement()
}

func (p *Parser) Function(kind string) Stmt {
	name := p.consume(lexer.IDENTIFIER, "Expect "+kind+" name.")
	p.consume(lexer.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	parameters := []lexer.Token{}
	if !p.check(lexer.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				r2loxerrors.Errors(p.peek().Line, "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(lexer.IDENTIFIER, "Expect parameter name."))
			if !p.match(lexer.COMMA) {
				break
			}
		}
	}
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(lexer.LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body := p.Block()
	return Function{Name: name, Parameters: parameters, Body: body}
}

func (p *Parser) VarDeclaration() Stmt {
	name := p.consume(lexer.IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(lexer.EQUAL) {
		initializer = p.Expression()
	}

	p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	return Var{Name: name, Initializer: initializer}
}

func (p *Parser) Statement() Stmt {
	if p.match(lexer.PRINT) {
		return p.PrintStatement()
	}

	if p.match(lexer.LEFT_BRACE) {
		return Block{Statements: p.Block()}
	}

	if p.match(lexer.IF) {
		return p.IfStatement()
	}

	if p.match(lexer.RETURN) {
		return p.ReturnStatement()
	}

	if p.match(lexer.WHILE) {
		return p.WhileStatement()
	}

	if p.match(lexer.FOR) {
		return p.ForStatement()
	}

	return p.ExpressionStatement()
}

func (p *Parser) Block() []Stmt {
	statements := []Stmt{}

	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.Declaration())
	}

	p.consume(lexer.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) IfStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.Expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.Statement()
	var elseBranch Stmt
	if p.match(lexer.ELSE) {
		elseBranch = p.Statement()
	}

	return If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

// while
func (p *Parser) WhileStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.Expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after while condition.")

	body := p.Statement()

	return While{Condition: condition, Body: body}
}

// for
func (p *Parser) ForStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer = p.VarDeclaration()
	} else {
		initializer = p.ExpressionStatement()
	}

	var condition Expr
	if !p.check(lexer.SEMICOLON) {
		condition = p.Expression()
	}
	p.consume(lexer.SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(lexer.RIGHT_PAREN) {
		increment = p.Expression()
	}
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.Statement()

	if increment != nil {
		body = Block{Statements: []Stmt{body, Expression{Expression: increment}}}
	}

	if condition == nil {
		condition = Literal{Value: true}
	}
	body = While{Condition: condition, Body: body}

	if initializer != nil {
		body = Block{Statements: []Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) Call() Expr {
	expr := p.Primary()

	for {
		if p.match(lexer.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}

	if !p.check(lexer.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				r2loxerrors.Errors(p.peek().Line, "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.Expression())
			if !p.match(lexer.COMMA) {
				break
			}
		}
	}

	paren := p.consume(lexer.RIGHT_PAREN, "Expect ')' after arguments.")

	return Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (p *Parser) PrintStatement() Stmt {
	value := p.Expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	return Print{Expression: value}
}

func (p *Parser) ExpressionStatement() Stmt {
	expr := p.Expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after expression.")
	return Expression{Expression: expr}
}

func (p *Parser) Expression() Expr {
	return p.assignment()
}

func (p *Parser) Parse() []Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.Declaration())
	}

	return statements

}

func (p *Parser) isAtEnd() bool {
	return p.Current >= len(p.Tokens)
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.Tokens[p.Current-1]
}
