package lexer

import (
	"testing"
	)

func TestNextToken(t *testing.T) {
	input :=`=+(){},;`

	tests := []struct{
		ExpectedType TokenType
		ExpectedValue string
	}{
		{ASSIGN,"="},
		{PLUS,"+"},
		{LPAREN,"("},
		{RPAREN,")"},
		{LBRACE,"{"},
		{RBRACE,"}"},
		{COMMA,","},
		{SEMICOLON,";"},
	}

	l := New(input)
	for i,test := range tests{
		tok := l.NextToken()

		t.Log("token:",i,tok)
		if tok.Type != test.ExpectedType{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, test.ExpectedType, tok.Type)
		}
		if tok.Value != test.ExpectedValue{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, test.ExpectedValue, tok.Type)
		}
	}
}

func TestLexer_NextToken(t *testing.T) {
	input := `
    let five = 5;
	let ten = 10;
	
	let add = fn(x,y){
		x + y
	};
	
	let result = add(five,ten);

    `

	tests := []struct{
		ExpectedType TokenType
		ExpectedValue string
	}{
		{LET,"let"},
		{INDENT,"five"},
		{ASSIGN,"="},
		{INT,"5"},
		{SEMICOLON,";"},

		{LET,"let"},
		{INDENT,"ten"},
		{ASSIGN,"="},
		{INT,"10"},
		{SEMICOLON,";"},

		{LET,"let"},
		{INDENT,"add"},
		{ASSIGN,"="},
		{FUNCTION,"fn"},
		{LPAREN,"("},
		{INDENT,"x"},
		{COMMA,","},
		{INDENT,"y"},
		{RPAREN,")"},

		{LBRACE,"{"},
		{INDENT,"x"},
		{PLUS,"+"},
		{INDENT,"y"},
		{RBRACE,"}"},
		{SEMICOLON,";"},

		{LET,"let"},
		{INDENT,"result"},
		{ASSIGN,"="},
		{INDENT,"add"},
		{LPAREN,"("},
		{INDENT,"five"},
		{COMMA,","},
		{INDENT,"ten"},
		{RPAREN,")"},
		{SEMICOLON,";"},
	}

	l := New(input)
	for i,test := range tests{
		tok := l.NextToken()

		t.Log("token:",i,tok)
		if tok.Type != test.ExpectedType{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, test.ExpectedType, tok.Type)
		}
		if tok.Value != test.ExpectedValue{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, test.ExpectedValue, tok.Type)
		}
	}
}