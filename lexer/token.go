package lexer

type TokenType string

const(
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	//标识符 + literature
	INDENT = "INDENT"
	INT = "INT"
	STRING = "STRING"

	//operator
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"

	LT = "<"
	GT = ">"

	EQ = "=="
	NOT_EQ = "!="

	//分隔符
	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	//KEYWORDS
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "true"
	FALSE = "false"
	IF = "if"
	ELSE = "else"
	RETURN = "return"

)

var keyWords = map[string]TokenType{
	"fn": FUNCTION,
	"let":LET,
	"true":TRUE,
	"false":FALSE,
	"if":IF,
	"else":ELSE,
	"return":RETURN,

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