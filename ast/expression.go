package ast

import (
		"lexer"
	"bytes"
		"strings"
)

type Expression interface {
	Node
	expressionNode()
}

type Indetifier struct {
	Token lexer.Token
	Value string
}

func (i *Indetifier)expressionNode(){}
func (i *Indetifier)TokenLiteral()string{
	return i.Token.Value
}
func (i *Indetifier)String()string{
	return i.Token.Value
}

type IntergerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntergerLiteral)expressionNode(){}

func (il *IntergerLiteral)TokenLiteral()string{
	return il.Token.Value
}
func (il *IntergerLiteral)String()string{
	return il.Token.Value
}

type PrefixExpression struct {
	Token lexer.Token //eg:!
	Operator string
	Right Expression
}

func (pe *PrefixExpression)expressionNode(){}

func (pe *PrefixExpression)TokenLiteral()string{
	return pe.Token.Value
}

func (pe *PrefixExpression)String()string{

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}


type InfixExpression struct {
	Token lexer.Token
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression)expressionNode(){}
func (ie *InfixExpression)TokenLiteral()string{
	return ie.Token.Value
}

func (ie *InfixExpression)String()string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(""+ie.Operator+"")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean)expressionNode(){}
func (b *Boolean)TokenLiteral()string{
	return b.Token.Value
}
func (b *Boolean)String()string{
	return b.Token.Value
}


type IfExpression struct {
	Token lexer.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression)expressionNode(){}
func (ie *IfExpression)TokenLiteral()string{
	return ie.Token.Value
}
func (ie *IfExpression)String()string{
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil{
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token lexer.Token //fn
	Parameters []*Indetifier
	Body *BlockStatement
}

func (fn *FunctionLiteral)expressionNode(){}
func (fn *FunctionLiteral)TokenLiteral()string{
	return fn.Token.Value
}
func (fn *FunctionLiteral)String()string{
	var out bytes.Buffer

	params := []string{}
	for _, p:= range fn.Parameters{
		params = append(params,p.String())
	}

	out.WriteString(fn.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params,","))
	out.WriteString(")")
	out.WriteString(fn.Body.String())

	return out.String()
}

type CallExpression struct {
	Token lexer.Token
	Function Expression
	Arguments []Expression
}

func (ce *CallExpression)expressionNode(){}
func (ce *CallExpression)TokenLiteral()string{
	return ce.Token.Value
}

func (ce *CallExpression)String()string{
	var out bytes.Buffer

	args := []string{}
	for _,a := range ce.Arguments{
		args = append(args,a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args,","))
	out.WriteString(")")

	return out.String()
}

//字符串常量
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (s *StringLiteral)expressionNode(){}
func (s *StringLiteral)TokenLiteral()string{
	return s.Token.Value
}
func (s *StringLiteral)String()string{
	return s.Token.Value
}

//数组
type ArrayLiteral struct {
	Token lexer.Token
	Element []Expression
}

func (al *ArrayLiteral)expressionNode()  {}
func (al *ArrayLiteral)TokenLiteral()string{
	return al.Token.Value
}
func (al *ArrayLiteral)String()string{
	var out bytes.Buffer

	eles := []string{}
	for _, el := range al.Element{
		eles = append(eles,el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(eles,","))
	out.WriteString("]")

	return out.String()
}


//数组下标
type IndexExpression struct {
	Token lexer.Token
	Left Expression
	Index Expression
}

func (ie *IndexExpression)expressionNode(){}
func (ie *IndexExpression)TokenLiteral()string{
	return ie.Token.Value
}
func (ie *IndexExpression)String()string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct {
	Token lexer.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral)expressionNode(){}
func (h *HashLiteral)TokenLiteral()string{
	return h.Token.Value
}
func (h *HashLiteral)String()string{
	var out bytes.Buffer

	pairs := []string{}
	for k,v := range h.Pairs{
		pairs = append(pairs,k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs,","))
	out.WriteString("}")

	return out.String()
}
