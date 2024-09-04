package parser

import (
	// "go/token"
	"fmt"

	"token/ast"
	lex "token/lexer"
	tok "token/token"
)

type Parser struct {
	l         *lex.Lexer // pointer to an instance of lexer which is called repeatedly(NextToken() ) to get the next token in input.
	curToken  tok.Token
	peekToken tok.Token
	errors    []string // initialize errors field

	prefixParseFns map[tok.TokenType]prefixParseFn
	infixParseFns map[tok.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

// stores the current and peek positions
func New(l *lex.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

// advances current and peek position
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != tok.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tok.LET:
		return p.parseLetStatement()
	case tok.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	p.nextToken()
	
	if !p.expectPeek(tok.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(tok.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(tok.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()


	for !p.curTokenIs(tok.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t tok.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tok.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tok.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}


//helper function for initializing errors
func (p *Parser) peekError(t tok.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}


func (p *Parser) registerPrefix(tokenType tok.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tok.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}