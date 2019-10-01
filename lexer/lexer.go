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

	l.skipSpace()

	switch l.char {
	case '*':
		tok = NewToken(ASTERISK,'*')
	case '!':
		if l.peekChar() == '='{
			tok = Token{NOT_EQ,"!="}
			l.readChar()
		}else{
			tok = NewToken(BANG,'!')
		}

	case '>':
		tok = NewToken(GT,'>')
	case '<':
		tok = NewToken(LT,'<')
	case '-':
		tok = NewToken(MINUS,'-')
	case '/':
		tok = NewToken(SLASH,'/')
	case '=':
		if l.peekChar() == '='{
			tok = Token{EQ,"=="}
			l.readChar()
		}else{
			tok =  NewToken(ASSIGN,'=')
		}

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
	case '"':// 字符串
		tok.Type = STRING
		tok.Value = l.readString()
	case '[':
		tok.Type = LBRACKET
		tok.Value = "["

	case ']':
		tok.Type = RBRACKET
		tok.Value = "]"
	case ':':
		tok.Type = COLON
		tok.Value = ":"

	default:
		if isLetter(l.char){
			tok.Value = l.readIdentifier() //这里已经预先读了一个字符，所以需要直接return
			tok.Type = LookIndentType(tok.Value)

			return tok
		}else if isDigit(l.char){
			tok.Value = l.readDigit()
			tok.Type = INT

			return tok
		}else{
			tok = NewToken(ILLEGAL,l.char)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer)peekChar()byte{
	if l.readPosition >= len(l.input){
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer)readIdentifier()string{
	position := l.position
	for isLetter(l.char){
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer)readDigit()string{
	position := l.position
	for isDigit(l.char){
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer)skipSpace(){
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r'{
		l.readChar()
	}
}


func isLetter(c byte)bool{
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c =='_'
}

func isDigit(c byte)bool{
	return c >= '0' && c <= '9'
}

func (l *Lexer)readString()string{
	pos := l.position + 1

	for {
		l.readChar()
		if l.char == '"'{
			break
		}
	}

	return l.input[pos:l.position]
}