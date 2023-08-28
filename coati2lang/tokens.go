package coati2lang

import (
	"fmt"
)

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN    TokenType = iota //[ok]
	RIGHT_PAREN                    //[ok]
	LEFT_BRACE                     //[ok]
	RIGHT_BRACE                    //[ok]
	LEFT_BRACKET                   //[]
	RIGHT_BRACKET                  //[]
	COMMA                          //[ok]
	DOT                            //[]
	MINUS                          //[ok]
	PLUS                           //[ok]
	SEMICOLON                      //[ok]
	SLASH                          //[ok]
	STAR                           //[ok]
	PERCENT                        //[]
	COLON                          //[]
	QUESTION                       //[]
	CARET                          //[]
	AMPERSAND                      //[]

	// One or two character tokens.
	BANG          //[ok]
	BANG_EQUAL    //[ok]
	EQUAL         //[ok]
	EQUAL_EQUAL   //[ok]
	GREATER       //[ok]
	GREATER_EQUAL //[ok]
	LESS          //[ok]
	LESS_EQUAL    //[ok]
	PLUS_PLUS     //[ok]
	MINUS_MINUS   //[ok]
	STAR_STAR     //[ok]
	ARROW         //[]
	LEFT          //[]
	RIGHT         //[]
	PIPE          //[OK] ABS
	OR_OR         //[]
	AND_AND       //[]

	// Literals.
	IDENTIFIER //[ok]
	STRING     //[ok]
	NUMBER     //[ok]

	// Keywords.
	AND        //[ok]
	CLASS      //[]
	ARRAY      //[]
	MAP        //[]
	MOD        //[]
	NOT        //[]
	TRY        //[]
	CATCH      //[]
	FINALLY    //[]
	THROW      //[]
	ADD        //[]
	DELETE     //[]
	TYPEOF     //[]
	INSTANCEOF //[]
	EXTENDS    //[]
	SWITCH     //[]
	CASE       //[]
	DEFAULT    //[]
	DO         //[]
	ELSE       //[ok]
	FALSE      //[ok]
	FUN        //[ok]
	IF         //[ok]
	NIL        //[ok]
	OR         //[ok]
	RETURN     //[ok]
	SUPER      //[]
	THIS       //[]
	TRUE       //[ok]
	VAR        //[ok]
	LET        //[]
	CONST      //[]
	WHILE      //[ok]
	FOR        //[ok]
	BREAK      //[]
	CONTINUE   //[]

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}

func ScanTokens(source string) []Token {
	// Implementación de tu escaneo de tokens aquí
	var tokens []Token
	scanner := NewScanner(source)
	tokens = scanner.ScanTokens()

	return tokens
}
