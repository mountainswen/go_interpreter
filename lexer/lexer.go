package lexer

type Lexer struct {
	input string //源码
	position int //position指向当前字符，即char所在的位置
	readPosition int //readposition指向下一个字符
	char byte
}

func New(input string)*Lexer{
	l := &Lexer{
		input:input,
	}

	l.readChar() //初始化
	return l
}

func (l *Lexer)readChar(){ //读取下一个字符
	if l.readPosition >= len(l.input) {
		l.char = 0
	}else{
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition = l.readPosition + 1
}

func (l *Lexer)NextToken()Token{
	var tok Token

	switch l.char {
	case '=':
		tok =  NewToken(ASSIGN,'=')
	case ';':
		tok =  NewToken(SEMICOLON,';')
	case '(':
		tok =  NewToken(LPAREN,'(')
	case ')':
		tok =  NewToken(RPAREN,')')
	case '{':
		tok =  NewToken(LBRACE,'{')
	case '}':
		tok =  NewToken(RBRACE,'}')
	case ',':
		tok =  NewToken(COMMA,',')
	case '+':
		tok =  NewToken(PLUS,'+')
	case 0:
		tok.Value = ""
		tok.Type = EOF
	}

	l.readChar()
	return tok
}