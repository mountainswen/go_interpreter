package lexer

import "testing"

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