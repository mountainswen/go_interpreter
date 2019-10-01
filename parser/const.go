package parser

import "lexer"

const (
	_ int = iota

	LOWEST
	EQUALS      // ==
	LESSGREATER // < >
	SUM         // +
	PRODUCT     //*
	PREFIX      //-X OR !X
	CALL        //fn(X)
	INDEX
)

var (
	precedence = map[lexer.TokenType]int{
		lexer.EQ:       EQUALS,
		lexer.NOT_EQ:   EQUALS,
		lexer.LT:       LESSGREATER,
		lexer.GT:       LESSGREATER,
		lexer.PLUS:     SUM,
		lexer.MINUS:    SUM,
		lexer.SLASH:    PRODUCT,
		lexer.ASTERISK: PRODUCT,
		lexer.LPAREN:CALL,
		lexer.LBRACKET:INDEX,
	}
)
