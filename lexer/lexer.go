package token

import "token/token"

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

func New(input string) *Lexer {
	l:=&Lexer{input: input} //l is assigned to a pointer to lexer struct input.
	l.readChar()
	return l
}
/*
readChar- gives the next character in the string and advance our position in the input string
if readposition is greater than len(input) ch is set to 0(equivalent to NUL in ASCII)


*/
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}


func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()//skip all forms of whitespace and new lines since they are only used as separators of tokens in Kisumu

	
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		}else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string{
	position := l.position
	for isLetter(l.ch){
		l.readChar()
	}
return l.input[position:l.position]
}

/*
NextToken -> Basically looks at the current character under examination l.ch and
returns the token depending on which character it is
newToken -> Helps in intializatio of the tokens
*/
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && 'z'<=ch || 'A' <= ch && 'Z' <= ch || ch == '_' || ch == '-'
}

//used for skipping whitespaces(also called eatWhitespace)
func (l *Lexer) skipWhitespace() {
for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
	l.readChar()
}
}

//check if its a digit
func isDigit(ch byte)bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}