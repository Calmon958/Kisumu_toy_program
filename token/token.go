package token

type TokenType string 

type Token struct{
	Type TokenType
	Literal string
}

const(
	//special types
	ILLEGAL = "ILLEGAL" //token or character not known about
	EOF = "EOF" //end of file

	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	MULTIPLY = "*"
	DIVIDE = "/"

	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET = "LET"
)
 
var keywords = map[string]TokenType {
	"fn" : FUNCTION,
	"let" : LET,
}

func LookUpIdent(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}