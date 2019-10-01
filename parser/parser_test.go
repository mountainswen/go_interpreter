package parser

import (
	"testing"
	"lexer"
	"ast"
	"fmt"
)

func TestParser_ParseProgram(t *testing.T) {
	input := `
		let x= 5;
		let y = 10;
		let foobar =  838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t,p)
	if program == nil{
		t.Fatal("parseProgram return nil")
	}

	if len(program.Statements) != 3{
		t.Fatalf("program.statement does not contains 3 statements,got=%d",
			len(program.Statements))
	}


	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i,tt := range tests{
		stmt := program.Statements[i]

		if !testLetStatement(t,stmt,tt.expectedIdentifier){
			return
		}
	}

}

func testLetStatement(t *testing.T,s ast.Statement,name string)bool{
	if s.TokenLiteral() != "let"{
		t.Errorf("s.tokenLiteral not let, got=%q",s.TokenLiteral())
		return false
	}

	letStmt,ok := s.(*ast.LetStatement)
	if !ok{
		t.Errorf("s not *asn.LetStatement,got=%T",s)
		return false
	}

	if letStmt.Name.Value != name{
		t.Errorf("letStmt.Name value not '%s'.got=%s",name,letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name{
		t.Errorf("s.name  not '%s'.got=%s",name,letStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T,p *Parser){
	errors := p.Errors()

	if len(errors) == 0{
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _,msg := range errors{
		t.Errorf("parse error:%q",msg)
	}

	t.FailNow()
}

func TestParser_ReturnStatement(t *testing.T){
	input := `
		return 5;
		return 10;
		return 818818;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 3{
		t.Fatalf("program statement does not contains 3 statement,instead:" +
			"%d", len(program.Statements))
	}

	for _,stmt := range program.Statements{
		returnStmt,ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.return stmt,got=%T",stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return"{
			t.Errorf("returnstatement.tokenliteral not 'return'" +
				"got %q",returnStmt.TokenLiteral())
		}
	}


}

func TestString(t *testing.T){
	program := &ast.Program{
		Statements:[]ast.Statement{
			&ast.LetStatement{
				Token:lexer.Token{
					Type:lexer.LET,
					Value:"let"},

				Name:&ast.Indetifier{
					Token:lexer.Token{Type:lexer.INDENT,Value:"myVar"},
					Value:"myValue",
				},

				Value:&ast.Indetifier{
					Token:lexer.Token{Type:lexer.INDENT,Value:"anotherVar"},
					Value:"anotherValue",
				},
			},
		},
	}

	if program.String() != "let myVar=anotherVar;"{
		t.Errorf("program.String wrong,got=%s",program.String())
	}
}

func TestIdentifierExpression(t *testing.T){
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1{
		t.Fatalf("program has not enough statements,got = %d",
			len(program.Statements))
	}

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("program.statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}

	ident,ok := stmt.Expression.(*ast.Indetifier)
	if !ok{
		t.Fatalf("exp not *ast.Identifier,got=%T",stmt.Expression)
	}

	if ident.Value != "foobar"{
		t.Errorf("ident.Value not %s,got=%s","foobar",
			ident.Value)
	}

	if ident.TokenLiteral() != "foobar"{
		t.Errorf("ident.TokenLiteral not %s, got=%s","foobar",
			ident.TokenLiteral())
	}
}

func TestIdentifierLiteralExpression(t *testing.T){
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1{
		t.Fatalf("program has not enough statements,got = %d",
			len(program.Statements))
	}

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("program.statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}

	interal,ok := stmt.Expression.(*ast.IntergerLiteral)
	if !ok{
		t.Fatalf("exp not *ast.IntergerLiteral,got=%T",stmt.Expression)
	}

	if interal.Value != 5{
		t.Errorf("ident.Value not %s,got=%d","5",
			interal.Value)
	}

	if interal.TokenLiteral() != "5"{
		t.Errorf("ident.TokenLiteral not %s, got=%s","5",
			interal.TokenLiteral())
	}
}

func TestParser_PrefixExpression(t *testing.T){
	prefixTests := []struct{
		input string
		operator string
		integerValue int64
	}{
		{"!5;","!",5},
		{"-15;","-",15},
	}

	for _,tt := range prefixTests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1{
			t.Fatalf("program.statements dose not contains %d statement" +
				" but got=%d",1, len(program.Statements))
		}

		stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("program.ststemt[0] is not ast.expressionStatement,got %T",
				program.Statements[0])
		}

		exp,ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmtm is not prefixExpression,got" +
				"=%T",stmt.Expression)
		}

		if exp.Operator != tt.operator{
			t.Fatalf("exp.Operator is not %s" +
				",got=%s",tt.operator,exp.Operator)
		}

		if !testIntegerLiteral(t,exp.Right,tt.integerValue){
			return
		}

	}
}

func testIntegerLiteral(t *testing.T,
	il ast.Expression,value int64)bool{
		integ,ok := il.(*ast.IntergerLiteral)
		if !ok{
			t.Errorf("il not ast.IntegerLiteral,got=%T",
			il)
			return false
		}

		if integ.Value != value{
			t.Errorf("integ.Value not %d,got = %d",
				value,integ.Value)
			return false
		}

		if integ.TokenLiteral() != fmt.Sprintf("%d",value){
			t.Errorf("integ.TokenLiteral not %d,got=%s",value,
				integ.TokenLiteral())
			return false
		}

		return true
}

func TestParse_InfixExpression(t *testing.T){
	infixTests := []struct{
		input string
		leftValue int64
		operator string
		rightValue int64
	}{
		{"5 + 5;",5 ,"+",5},
		{"5 - 5;",5,"-",5},
		{"5 * 5;",5,"*",5},
		{"5/5;",5,"/",5},
		{"5 > 5;",5,">",5},
		{"5 < 5;",5,"<",5},
		{"5 == 5;",5,"==",5},
		{"5 != 5",5,"!=",5},
	}

	for _,tt := range infixTests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1{
			t.Fatalf("program.Statements does not " +
				"contains %d statements,got=%d\n",1,
				len(program.Statements))
		}

		stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("program.Statements[0] is not st.expressionStatement" +
				" ,got=%T",program.Statements[0])
		}

		exp,ok := stmt.Expression.(*ast.InfixExpression)
		if !ok{
			t.Fatalf(" exp is not infix expression" +
				" but got=%T ",exp)
		}

		if !testIntegerLiteral(t,exp.Left,tt.leftValue){
			return
		}

		if exp.Operator != tt.operator{
			t.Fatalf("exp.Operator is not %s,got=" +
				"%s",tt.operator,exp.Operator)
		}

		if !testIntegerLiteral(t,exp.Right,tt.rightValue){
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T){
	tests := []struct{
		input string
		expected string
	}{
		{"-a * b","((-a)*b)"},
		{"!-a","(!(-a))"},
		{"a + b + c", "((a+b)+c)"},
		{"a+b-c","((a+b)-c)"},
		{"a*b*c","((a*b)*c)"},
		{"a * b/c","((a*b)/c)"},
		{"a+b/c","(a+(b/c))"},
		{"a+b*c+d/e-f","(((a+(b*c))+(d/e))-f)"},
		{"3+4;-5*5","(3+4)((-5)*5)"},
		{"5>4 == 3<4","((5>4)==(3<4))"},
		{"5<4 != 3 > 4","((5<4)!=(3>4))"},
		{"3 + 4* 5 == 3 * 1 + 4* 5","((3+(4*5))==((3*1)+(4*5)))" },
	}

	for _,tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t,p)

		actual := program.String()

		if actual != tt.expected{
			t.Errorf("expected=%q,got=%q",tt.expected,actual)
		}
	}
}

func TestParser_GroupExpression(t *testing.T){
	tests := []struct{
		input string
		expected string
	}{
		{"1+(2+3)","(1+(2+3))"},
		{"3*(4+5)","(3*(4+5))"},
	}

	for _,tt := range tests{
		input := tt.input
		l := lexer.New(input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1{
			t.Fatalf("expected %d statement but " +
				" got =%d",1, len(program.Statements))
		}

		actual := program.String()
		if actual != tt.expected{
			t.Fatalf("expected %s,but got=%s",
				tt.expected,actual)
		}

	}
}

/*
func TestIfExpression(t *testing.T) {
	input := `if(x<y){
					x
                }else{y}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1{
		t.Fatalf("program.Body does not contains " +
			" %d statements but got=%d",1, len(program.Statements))
	}


	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("program.statement[0] is not ast.expression" +
			" got=%T",stmt)
	}

	exp,ok := stmt.Expression.(*ast.IfExpression)
	if !ok{
		t.Fatalf("stmt.Expression is not ast.IfExpression" +
			", got =%T",exp)
	}

	//t.Error("ifexpression:",exp.Condition.String())
	//t.Error("if expression consequence:",exp.Consequence.String())
	//t.Error("ifexpression alternative",exp.Alternative.String())

}
*/
func TestParser_FunctionLiteral(t *testing.T){
	tests := []struct{
		input string
	}{
		{"fn(){return x,y}"},
		//{"fan(){}"},
		//{"fan(x,y){}"},
	}

	for _,tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1{
			t.Fatalf("expected %d statement,but got=%d",1, len(program.Statements))
		}

		stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("expected ast.expressionStatement,but got=%T",stmt)
		}

		exp,ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok{
			t.Fatalf("expected functional expression,but got=%T",stmt.Expression)
		}

		t.Errorf(exp.String())
	}
}
