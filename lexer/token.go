package lexer

type TokenType string

const(
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	//标识符 + literature
	INDENT = "INDENT"
	INT = "INT"

	//operator
	ASSIGN = "="
	PLUS = "+"

	//分隔符
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//KEYWORDS
	FUNCTION = "FUNCTION"
	LET = "LET"
)

var keyWords = map[string]TokenType{
	"fn": FUNCTION,
	"let":LET,
}

type Token struct {
	Type TokenType
	Value string
}

func NewToken(Type TokenType,Value byte) Token {
	return Token{Type:Type,Value:string(Value)}
}

func LookIndentType(indent string)TokenType{
	if tokType,ok := keyWords[indent];ok{
		return tokType
	}

	return INDENT
}