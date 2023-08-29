package coati2lang

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Parser struct {
	Tokens  []Token
	Current int
	Start   int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
		Start:   0,
	}
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) ReturnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.Expression()
	}

	p.consume(SEMICOLON, "Expect ';' after return value.")
	return Return{Keyword: keyword, Value: value}
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	for p.match(MINUS_MINUS) {
		op := p.previous()
		operator := Token{Type: MINUS, Lexeme: "--", Literal: nil, Line: op.Line}
		right := Literal{Value: 1.0}
		var name string
		if expr == nil {
			expr = p.Factor()
		}
		if expr, ok := expr.(Var); ok {
			name = expr.Name.Lexeme
		}

		expr = Assign{Name: Token{
			Type:    IDENTIFIER,
			Lexeme:  name,
			Literal: nil,
			Line:    op.Line,
		}, Value: Binary{Left: expr, Operator: operator, Right: right}}

	}

	for p.match(PLUS_PLUS) {
		op := p.previous()
		operator := Token{Type: PLUS, Lexeme: "++", Literal: nil, Line: op.Line}
		right := Literal{Value: 1.0}
		var name string
		if expr == nil {
			expr = p.Factor()
		}
		if expr, ok := expr.(Var); ok {
			name = expr.Name.Lexeme
		}

		expr = Assign{Name: Token{
			Type:    IDENTIFIER,
			Lexeme:  name,
			Literal: nil,
			Line:    op.Line,
		}, Value: Binary{Left: expr, Operator: operator, Right: right}}
	}

	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()

	for p.match(SLASH, STAR, STAR_STAR, PERCENT) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		value := p.Unary()
		return Unary{Operator: operator, Value: value}
	}

	return p.Call()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{Value: false}
	}
	if p.match(TRUE) {
		return Literal{Value: true}
	}
	if p.match(NIL) {
		return Literal{Value: nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{Value: p.previous().Literal}
	}

	if p.match(IDENTIFIER) {
		name := p.previous()

		selectors := [][]Expr{}
		for p.check(LEFT_BRACKET) || p.check(DOT) {
			if p.match(LEFT_BRACKET) {
				array := p.Array()
				selectors = append(selectors, array)
			} else if p.match(DOT) {
				name := p.consume(IDENTIFIER, "Expect property name after '.'.")
				selectors = append(selectors, []Expr{Literal{Value: name.Lexeme}})
			}
		}

		return Var{Name: name, Selectors: selectors}
	}

	if p.match(LEFT_PAREN) {
		expr := p.Equality()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{Expression: expr}
	}

	if p.match(LEFT_BRACKET) {
		array := p.Array()
		return Literal{Value: array}
	}

	if p.match(LEFT_BRACE) {
		array := p.Map()
		return Literal{Value: array}
	}

	if p.match(PIPE) {
		expr := p.Equality()
		p.consume(PIPE, "Expect '|' after expression.")
		return GroupingABS{Expression: expr}
	}

	if p.match(EOF) {
		return Literal{Value: nil}
	}

	return nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) consume(t TokenType, message string) Token {
	line := p.peek().Line
	if p.check(t) {
		return p.advance()
	}

	message1 := fmt.Sprint(message, "\n\tLine -> ", fmt.Sprint(line), "\n\tNear ->"+p.peek().Lexeme)
	panic(message1)
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr, ok := expr.(Var); ok {
			name := expr.Name

			return Assign{Name: name, Value: value, Selectors: expr.Selectors}
		}

		Errors(equals.Line, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.Equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.Equality()
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) Declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%s", r)
			os.Exit(2)
			//p.synchronize()
		}
	}()

	if p.match(FUN) {
		name := p.consume(IDENTIFIER, "Expect function name.")
		return p.Function("function", name)
	}

	if p.match(VAR) {
		return p.VarDeclaration()
	}

	if p.match(LET) {
		return p.VarDeclaration()
	}

	if p.match(EOF) {
		return nil
	}
	return p.Statement()
}

func (p *Parser) Function(kind string, name Token) Stmt {

	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")

	parameters := []Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				Errors(p.peek().Line, "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body := p.Block()
	return Function{Name: name, Parameters: parameters, Body: body, Closure: nil}
}

func (p *Parser) VarDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	if p.match(EQUAL) {
		initializer := p.Expression()
		p.consume(SEMICOLON, "Expect ';' after variable declaration.")
		return Var{Name: name, InitializerVal: initializer}
	}
	if p.match(LEFT_BRACKET) {
		inizializer := []Expr{}
		if p.check(NUMBER) {
			size := p.consume(NUMBER, "Expect size of array.")
			size_declarate := int(size.Literal.(float64))
			p.consume(RIGHT_BRACKET, "Expect ']' after arguments.")
			if p.check(EQUAL) {
				p.consume(EQUAL, "Expect '=' after ']'.")
				p.consume(LEFT_BRACKET, "Expect '[' after '='.")
				inizializer = p.Array()

				if size_declarate != len(inizializer) {
					Errors(size.Line, "Size of array and number of arguments must be the same.")
				}
			}
			p.consume(SEMICOLON, "Expect ';' after variable declaration.")

			return Var{Name: name, InitializerArray: inizializer, SizeArrayInit: size_declarate}
		}

		p.consume(RIGHT_BRACKET, "Expect ']' after arguments.")
		if p.check(EQUAL) {
			p.consume(EQUAL, "Expect '=' after ']'.")
			p.consume(LEFT_BRACKET, "Expect '[' after '='.")
			initializer := p.Array()
			p.consume(SEMICOLON, "Expect ';' after variable declaration.")
			return Var{Name: name, InitializerArray: initializer, SizeArrayInit: len(initializer)}
		}

		p.consume(SEMICOLON, "Expect ';' after variable declaration.")
		return Var{Name: name, InitializerArray: inizializer, SizeArrayInit: 0}
	}

	if p.match(LEFT_BRACE) {
		p.consume(RIGHT_BRACE, "Expect '}' after arguments..")
		p.consume(EQUAL, "Expect '=' after '}'.")
		p.consume(LEFT_BRACE, "Expect '{' after '='.")
		initializer := p.Map()
		p.consume(SEMICOLON, "Expect ';' after variable declaration.")
		return Var{Name: name, InitializerMap: initializer}
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return Var{Name: name, InitializerVal: nil}
}

func (p *Parser) Array() []Expr {
	initializer := []Expr{}

	if !p.check(RIGHT_BRACKET) {
		for {
			if len(initializer) >= 255 {
				Errors(p.peek().Line, "Can't have more than 255 arguments.")
			}
			expr := p.Expression()

			if p.match(DOT) {
				p.consume(DOT, "Expect '.' after arguments.")
				literal := p.Expression()
				start := int(expr.(Literal).Value.(float64))
				end := int(literal.(Literal).Value.(float64))
				for i := start; i <= end; i++ {
					initializer = append(initializer, Literal{Value: float64(i)})
				}
			} else if p.match(LEFT_BRACKET) {
				subarray := p.Array()
				subVar := Var{Name: Token{Type: IDENTIFIER, Lexeme: "subarray-" + uuid.NewString()}, Sub: true, InitializerArray: subarray}
				initializer = append(initializer, subVar)
			} else if p.match(LEFT_BRACE) {
				submap := p.Map()
				subVar := Var{Name: Token{Type: IDENTIFIER, Lexeme: "submap-" + uuid.NewString()}, Sub: true, InitializerMap: submap}
				initializer = append(initializer, subVar)
			} else if p.match(FUN) {
				name := Token{Type: IDENTIFIER, Lexeme: "subfx-" + uuid.NewString()}
				subfx := p.Function("method", name)
				subVar := Var{Name: name, Sub: true, InitializerFx: subfx.(Function)}
				initializer = append(initializer, subVar)
			} else {
				initializer = append(initializer, expr)
			}

			if !p.match(COMMA) {
				break
			}
		}
	}

	p.consume(RIGHT_BRACKET, "Expect ']' after arguments.")

	return initializer
}

func (p *Parser) Map() []ItemVar {
	initializer := []ItemVar{}

	if !p.check(RIGHT_BRACE) {
		for {

			if len(initializer) >= 255 {
				Errors(p.peek().Line, "Can't have more than 255 arguments.")
			}
			key := p.Expression()

			if key_identifier, ok := key.(Var); ok {
				key = Literal{Value: key_identifier.Name.Lexeme}
			}

			if p.check(ARROW) {
				p.consume(ARROW, "Expect '=>' after key.")
				name := Token{Type: IDENTIFIER, Lexeme: "subfx-" + uuid.NewString()}
				subfx := p.Function("method", name)
				subVar := Var{Name: name, Sub: true, InitializerFx: subfx.(Function)}
				initializer = append(initializer, ItemVar{Key: key, Value: subVar})

			} else {

				if p.check(RIGHT_BRACE) {
					break
				}
				p.consume(COLON, "Expect ':' after key.")

				if p.match(LEFT_BRACKET) {
					subarray := p.Array()
					subVar := Var{Name: Token{Type: IDENTIFIER, Lexeme: "subarray-" + uuid.NewString()}, Sub: true, InitializerArray: subarray}
					initializer = append(initializer, ItemVar{Key: key, Value: subVar})
				} else if p.match(LEFT_BRACE) {
					submap := p.Map()
					subVar := Var{Name: Token{Type: IDENTIFIER, Lexeme: "submap-" + uuid.NewString()}, Sub: true, InitializerMap: submap}
					initializer = append(initializer, ItemVar{Key: key, Value: subVar})
				} else if p.match(FUN) {
					name := Token{Type: IDENTIFIER, Lexeme: "subfx-" + uuid.NewString()}
					subfx := p.Function("method", name)
					subVar := Var{Name: name, Sub: true, InitializerFx: subfx.(Function)}
					initializer = append(initializer, ItemVar{Key: key, Value: subVar})
				} else {
					value := p.Expression()
					initializer = append(initializer, ItemVar{Key: key, Value: value})
				}

			}

			if p.check(RIGHT_BRACE) {
				break
			}
			if p.check(SEMICOLON) {
				p.consume(SEMICOLON, "Expect ';' after arguments.")
				continue
			}

			if p.check(COMMA) {
				p.consume(COMMA, "Expect ',' after value.")
				continue
			}

		}
	}

	p.consume(RIGHT_BRACE, "Expect '}' after arguments.")
	return initializer
}

func (p *Parser) Statement() Stmt {
	if p.match(LEFT_BRACE) {
		return Block{Statements: p.Block()}
	}

	if p.match(IF) {
		return p.IfStatement()
	}

	if p.match(RETURN) {
		return p.ReturnStatement()
	}

	if p.match(WHILE) {
		return p.WhileStatement()
	}

	if p.match(FOR) {
		return p.ForStatement()
	}

	return p.ExpressionStatement()
}

func (p *Parser) Block() []Stmt {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.Declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) IfStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.Expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.Statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.Statement()
	}

	return If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) WhileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.Expression()
	p.consume(RIGHT_PAREN, "Expect ')' after while condition.")
	body := p.Statement()
	return While{Condition: condition, Body: body}
}

func (p *Parser) ForStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.VarDeclaration()
	} else {
		initializer = p.ExpressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.Expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.Expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

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
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}

	var this Var
	if self, ok := callee.(Var); ok {
		if len(self.Selectors) > 0 {
			this = self
			this.Selectors = this.Selectors[:len(this.Selectors)-1]
		}
	}

	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				Errors(p.peek().Line, "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.Expression())
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return Call{Callee: callee, Paren: paren, Arguments: arguments, This: this}
}

func (p *Parser) ExpressionStatement() Stmt {
	var stmt Stmt

	expr := p.Expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	stmt = Expression{Expression: expr}

	return stmt
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

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.Tokens[p.Current-1]
}
