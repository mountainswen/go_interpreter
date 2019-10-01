package parser

import (
	"ast"
	"lexer"
)

type (
	prefixParsefn func() ast.Expression
	infoxParsefn func(ast.Expression) ast.Expression
)


func (p *Parser)registerPrefix(tokenType lexer.TokenType,fn prefixParsefn){
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser)registerInfix(tokenType lexer.TokenType,fn infoxParsefn){
	p.infixParseFns[tokenType] = fn
}