package token

type TokenType string 

type Token struct{
	Type TokenType
	Literal string
}

const(
//essential strings are keyword tokens. Anything assigned to a string is a key word character.
// Those assigned to a single character are known as one-character tokens
// tw0-character tokens ie '==', '!='

	//special types
	ILLEGAL = "ILLEGAL" //token or character not known about
	EOF = "EOF" //end of file

	IDENT = "IDENT"
	INT = "INT"

	//operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	ASTERIC = "*"
	SLASH = "/"
	BANG = "!"

	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="

	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"


	//keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
	// WHILE = "WHILE"
	TRUE = "TRUE"
	FALSE = "FALSE"

	//two-character tokens
	EQ = "=="
	NOT_EQ = "!="
)
 
var keywords = map[string]TokenType {
	"fn" : FUNCTION,
	"let" : LET,
	"true" : TRUE,
	"false" : FALSE,
	"if" : IF,
	"else" : ELSE,
	"return" : RETURN,
	// "while" : WHILE,

}

func LookUpIdent(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}