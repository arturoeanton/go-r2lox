package lexer

import (
	"strconv"

	"github.com/arturoeanton/go-r2lox/r2loxerrors"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
	Line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source:  source,
		Tokens:  []Token{},
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// Estamos al comienzo del siguiente lexema.
		s.Start = s.Current
		s.scanToken() // Asumiendo que 'scanToken' está definido y toma estos argumentos
	}

	s.Tokens = append(s.Tokens, NewToken(EOF, "", nil, s.Line))
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, "(")
	case ')':
		s.addToken(RIGHT_PAREN, ")")
	case '{':
		s.addToken(LEFT_BRACE, "{")
	case '}':
		s.addToken(RIGHT_BRACE, "}")
	case ',':
		s.addToken(COMMA, ",")
	case '.':
		s.addToken(DOT, ".")
	case '-':
		s.addToken(MINUS, "-")
	case '+':
		s.addToken(PLUS, "+")
	case ';':
		s.addToken(SEMICOLON, ";")
	case '*':
		s.addToken(STAR, "*")
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, "!=")
		} else {
			s.addToken(BANG, "!")
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, "==")
		} else {
			s.addToken(EQUAL, "=")
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, ">=")
		} else {
			s.addToken(LESS, "<")
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, ">=")
		} else {
			s.addToken(GREATER, ">")
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, "/")
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.Line++
	case '"':
		s.string()

	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			r2loxerrors.Errors(s.Line, "Unexpected character.")
		}
	}
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	if err != nil {
		r2loxerrors.Errors(s.Line, "Error parsing number.")
	}

	s.addToken(NUMBER, value)
}

func (s *Scanner) peekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return '\x00'
	}
	return rune(s.Source[s.Current+1])
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// Verifica si la palabra es una palabra reservada.
	text := s.Source[s.Start:s.Current]

	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType, text)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00' // Carácter nulo en Go
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}
	s.advance()

	if s.isAtEnd() {
		r2loxerrors.Errors(s.Line, "Unterminated string.")
		return
	}

	value := s.Source[s.Start+1 : s.Current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.Source[s.Current]) != expected {
		return false
	}
	s.Current++
	return true
}

func (s *Scanner) advance() rune {
	c := rune(s.Source[s.Current])
	s.Current++
	return c
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.Source[s.Start:s.Current]
	token := NewToken(tokenType, text, literal, s.Line)
	s.Tokens = append(s.Tokens, token)
}
