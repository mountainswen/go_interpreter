package ast

import (
	"lexer"
	"bytes"
)

type Statement interface {
	Node
	statmentNode()
}

type LetStatement struct {
	Token lexer.Token
	Name *Indetifier
	Value Expression
}

func (l *LetStatement)statmentNode(){}
func (l *LetStatement)TokenLiteral()string{
	return l.Token.Value
}
func (l *LetStatement)String()string{
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString("=")

	if l.Value != nil{
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement)statmentNode(){}
func (rs *ReturnStatement)TokenLiteral()string{
	return rs.Token.Value
}

func (rs *ReturnStatement)String()string{
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil{
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}


type ExpressionStatement struct {
	Token lexer.Token
	Expression Expression
}

func (e *ExpressionStatement)statmentNode(){}
func (e *ExpressionStatement)TokenLiteral()string{
	return e.Token.Value
}

func (e *ExpressionStatement)String()string{
	if e.Expression != nil{
		return e.Expression.String()
	}

	return ""
}

type BlockStatement struct {
	Token lexer.Token
	Statements []Statement
}

func (b *BlockStatement)String()string{
	var out bytes.Buffer

	for _,s := range b.Statements{
		out.WriteString(s.String())
	}

	return out.String()
}

func (b *BlockStatement)TokenLiteral()string{
	return b.Token.Value
}
func (b *BlockStatement)statmentNode(){

}
