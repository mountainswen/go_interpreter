package parser

import "lexer"
import (
	"ast"
	"fmt"
	"strconv"
	)

//parser的作用是将源码翻译成一个数据结构
/*
“A parser is a software component that takes input data (frequently text) and
	builds a data structure”

摘录来自: Ball, Thorsten. “Writing An Interpreter In Go。” iBooks.
*/
type Parser struct {
	l *lexer.Lexer

	errors []string
	curToken lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParsefn //是个tokentype
	infixParseFns map[lexer.TokenType]infoxParsefn
}

func New(l *lexer.Lexer)*Parser{
	p := &Parser{l:l,errors:[]string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenType]prefixParsefn)
	p.registerPrefix(lexer.INDENT,p.parseIdentifier)
	p.registerPrefix(lexer.INT,p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING,p.parseStringLiteral)

	p.registerPrefix(lexer.MINUS,p.parsePrefixExpression)
	p.registerPrefix(lexer.BANG,p.parsePrefixExpression)

	p.registerPrefix(lexer.TRUE,p.parseBoolean)
	p.registerPrefix(lexer.FALSE,p.parseBoolean)

	p.registerPrefix(lexer.LPAREN,p.parseGroupExpression)
	p.registerPrefix(lexer.IF,p.parseIfExpression)
	p.registerPrefix(lexer.FUNCTION,p.parseFunctionLiteral)
	p.registerPrefix(lexer.LBRACKET,p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE,p.parseHashLiteral)
	//infix
	p.infixParseFns = make(map[lexer.TokenType]infoxParsefn)
	p.registerInfix(lexer.PLUS,p.parseInfixExpression)
	p.registerInfix(lexer.MINUS,p.parseInfixExpression)
	p.registerInfix(lexer.SLASH,p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK,p.parseInfixExpression)
	p.registerInfix(lexer.EQ,p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ,p.parseInfixExpression)
	p.registerInfix(lexer.LT,p.parseInfixExpression)
	p.registerInfix(lexer.GT,p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN,p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET,p.parseIndexExpression)
	return p
}

func (p *Parser)parseHashLiteral()ast.Expression{
	hash := &ast.HashLiteral{Token:p.curToken}

	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenis(lexer.RBRACE){
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.COLON){
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenis(lexer.RBRACE) && !p.expectPeek(lexer.COMMA){
			return nil
		}
	}

	if !p.expectPeek(lexer.RBRACE){
		return nil
	}

	return hash
}

func (p *Parser)parseIndexExpression(array ast.Expression)ast.Expression{
	idxExp := &ast.IndexExpression{}
	idxExp.Left = array

	p.nextToken()
	idxExp.Index = p.parseExpression(LOWEST)

	if !p.peekTokenis(lexer.RBRACKET){
		return nil
	}

	p.nextToken()
	return idxExp
}

func (p *Parser)parseArrayLiteral()ast.Expression{
	array := &ast.ArrayLiteral{}

	array.Element = p.parseExpressionList(lexer.RBRACKET)

	return array
}

func (p *Parser)parseExpressionList(end lexer.TokenType)[]ast.Expression{
	list := []ast.Expression{}

	if p.peekTokenis(end){
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list,p.parseExpression(LOWEST))

	for p.peekTokenis(lexer.COMMA){
		p.nextToken()
		p.nextToken()
		list = append(list,p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end){
		return nil
	}

	return list
}

func (p *Parser)Errors()[]string{
	return p.errors
}

func (p *Parser)peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s,got %s instead", t, p.peekToken.Type)

	p.errors = append(p.errors,msg)
}

func (p *Parser)nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 明天实现一下这个函数
func (p *Parser)ParseProgram()*ast.Program{

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF{
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements,stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser)ParseStatement()ast.Statement{
	switch p.curToken.Type {
	case lexer.LET:
		return p.ParseLetStatement()
	case lexer.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.parseExpressionStatement()
	//	msg := fmt.Sprintf("invalid statement")
	//	p.errors = append(p.errors,msg)

	}
	return nil
}

func (p *Parser)ParseReturnStatement()*ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token:p.curToken}

	p.nextToken()

	if !p.curTokenis(lexer.SEMICOLON){
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	p.nextToken()

	return stmt
}

func (p *Parser)ParseLetStatement()*ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(lexer.INDENT) {
		return nil
	}

	stmt.Name = &ast.Indetifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.expectPeek(lexer.ASSIGN){
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	p.nextToken()
	for p.curTokenis(lexer.SEMICOLON){
		p.nextToken()
	}


	return stmt

}

func (p *Parser)curTokenis(t lexer.TokenType)bool{
	return p.curToken.Type == t
}

func (p *Parser)peekTokenis(t lexer.TokenType)bool{
	return p.peekToken.Type == t
}

func (p *Parser)expectPeek(t lexer.TokenType)bool{
	if p.peekTokenis(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}

}


func (p *Parser)parseExpressionStatement()*ast.ExpressionStatement{
	stmt := &ast.ExpressionStatement{}
	stmt.Token = p.curToken

	stmt.Expression = p.parseExpression(LOWEST) //最低优先级

	if p.peekTokenis(lexer.SEMICOLON){
		p.nextToken()
	}

	return stmt
}

func (p *Parser)parseExpression(precedence int)ast.Expression{
	curTok := p.curToken

	prefix := p.prefixParseFns[curTok.Type]

	if prefix == nil{
		p.noPrefixParseFnError(p.curToken)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenis(lexer.SEMICOLON) &&
		precedence < p.peekPrecedence(){
			infix := p.infixParseFns[p.peekToken.Type]
			if infix == nil{
				return leftExp
			}

			p.nextToken()

			leftExp = infix(leftExp)
	}

	return leftExp

}

func (p *Parser)parseIdentifier()ast.Expression{
	return &ast.Indetifier {
		Token:p.curToken,
		Value:p.curToken.Value,
	}
}

func (p *Parser)parseIntegerLiteral()ast.Expression{
	lit := &ast.IntergerLiteral{Token:p.curToken}

	value,err := strconv.ParseInt(p.curToken.Value,10,64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer",
			p.curToken.Value)

		p.errors = append(p.errors,msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser)parseStringLiteral()ast.Expression{
	lit := &ast.StringLiteral{}

	lit.Value = p.curToken.Value

	return lit
}

func (p *Parser)parsePrefixExpression()ast.Expression{
	expression := &ast.PrefixExpression{
		Token:p.curToken,
		Operator:p.curToken.Value,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser)parseInfixExpression(left ast.Expression)ast.Expression{
	stmt := &ast.InfixExpression{
		Token:p.curToken,
		Operator:p.curToken.Value,
		Left:left,
	}

	precedence := p.curPrecedence()

	p.nextToken()

	stmt.Right = p.parseExpression(precedence)
	return stmt
}

func (p *Parser)noPrefixParseFnError(t lexer.Token){
	msg := fmt.Sprintf("no prefix parse function for" +
		" %s found ",t.Type)

	p.errors = append(p.errors,msg)
}

func (p *Parser)peekPrecedence()int{
	if p,ok := precedence[p.peekToken.Type];ok{
		return p
	}

	return LOWEST
}

func (p *Parser)curPrecedence()int{
	if p,ok := precedence[p.curToken.Type];ok{
		return p
	}

	return LOWEST
}

func (p *Parser)parseBoolean()ast.Expression{
	return &ast.Boolean{
		Token:p.curToken,
		Value:p.curTokenis(lexer.TRUE),
	}
}

func (p *Parser)parseGroupExpression()ast.Expression{
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN){
		return nil
	}

	return exp
}

func (p *Parser)parseIfExpression()ast.Expression{
	expression := &ast.IfExpression{}

	if !p.expectPeek(lexer.LPAREN){
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN){
		return nil
	}

	if !p.expectPeek(lexer.LBRACE){
		return nil
	}

	expression.Consequence =p.parseBlockStatement()

	if p.peekTokenis(lexer.ELSE){
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE){
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser)parseBlockStatement()*ast.BlockStatement{
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenis(lexer.RBRACE) && !p.curTokenis(lexer.EOF){
		stmt := p.ParseStatement()
		if stmt != nil{
			block.Statements = append(block.Statements,stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser)parseFunctionLiteral()ast.Expression{

	function := &ast.FunctionLiteral{
		Token:p.curToken,
	}

	if !p.expectPeek(lexer.LPAREN){
		return nil
	}

	function.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.LBRACE){
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser)parseFunctionParameters()[]*ast.Indetifier{

	identifiers := []*ast.Indetifier{}

	if p.peekTokenis(lexer.RPAREN){
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Indetifier{Token:p.curToken,Value:p.curToken.Value}

	identifiers = append(identifiers,ident)
	for p.peekTokenis(lexer.COMMA){
		p.nextToken()
		p.nextToken()

		ident := &ast.Indetifier{Token:p.curToken,Value:p.curToken.Value}
		identifiers = append(identifiers,ident)
	}

	if !p.expectPeek(lexer.RPAREN){
		return nil
	}

	return identifiers
}

func (p *Parser)parseCallExpression(function ast.Expression)ast.Expression{
	exp := &ast.CallExpression{
		Token:p.curToken,
		Function:function,
	}

	exp.Arguments = p.parseCallArgument()

	return exp
}

func (p *Parser)parseCallArgument()[]ast.Expression{
	args := []ast.Expression{}

	if p.peekTokenis(lexer.RPAREN){
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args,p.parseExpression(LOWEST))

	for p.peekTokenis(lexer.COMMA){
		p.nextToken()
		p.nextToken()

		args = append(args,p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RPAREN){
		return nil
	}

	return args
}