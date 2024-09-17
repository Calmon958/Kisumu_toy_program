package parser

import (
	// "go/token"
	"fmt"
	"strconv"

	"token/ast"
	lex "token/lexer"
	// "token/token"

	// "token/token"
	tok "token/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -S or !S
	CALL        // add(a, b)
)

type Parser struct {
	l         *lex.Lexer // pointer to an instance of lexer which is called repeatedly(NextToken() ) to get the next token in input.
	curToken  tok.Token
	peekToken tok.Token
	errors    []string // initialize errors field

	prefixParseFns map[tok.TokenType]prefixParseFn
	infixParseFns  map[tok.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// stores the current and peek positions
func New(l *lex.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	// for the identifiers
	p.prefixParseFns = make(map[tok.TokenType]prefixParseFn)
	p.registerPrefix(tok.IDENT, p.parserIdentifier)

	// for integer
	p.prefixParseFns = make(map[tok.TokenType]prefixParseFn)
	p.registerPrefix(tok.IDENT, p.parserIdentifier)
	p.registerPrefix(tok.INT, p.parseIntegerLiteral)

	// prefix expressions and registers
	p.registerPrefix(tok.BANG, p.parserPrefixExpression)
	p.registerPrefix(tok.MINUS, p.parserPrefixExpression)

	//infix expression
	p.infixParseFns = make(map[tok.TokenType]infixParseFn)
	p.registerInfix(tok.PLUS, p.parseInfixExpression)
	p.registerInfix(tok.MINUS, p.parseInfixExpression)
	p.registerInfix(tok.SLASH, p.parseInfixExpression)
	p.registerInfix(tok.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tok.EQ, p.parseInfixExpression)
	p.registerInfix(tok.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tok.LT, p.parseInfixExpression)
	p.registerInfix(tok.GT, p.parseInfixExpression)

	//Boolean
	p.registerPrefix(tok.TRUE, p.parseBoolean)
	p.registerPrefix(tok.FALSE, p.parseBoolean)

	//Grouped
	p.registerPrefix(tok.LPAREN, p.parseGroupedExpression)

	//ifExpression
	p.registerPrefix(tok.IF, p.parseIfExpression)

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
		return p.parseExpressionStatement()
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

// helper function for initializing errors
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(tok.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t tok.TokenType) {
	msg := fmt.Sprintf("no prefix function for %s found", t)
	p.errors = append(p.errors, msg)
}

// checks whether there is a parser function associated with p.curToken.Type prefix
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(tok.SEMICOLON) && precedence < p. peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parserIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parserPrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parserExpression(PREFIX)
	return expression
}


var precedences = map[tok.TokenType]int {
	tok.EQ: EQUALS,
	tok.NOT_EQ: EQUALS,
	tok.LT: LESSGREATER,
	tok.GT: LESSGREATER,
	tok.PLUS: SUM,
	tok.MINUS: SUM,
	tok.SLASH: PRODUCT,
	tok.ASTERISK: PRODUCT,
}

//precedences associate token types with their precedences
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}



func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression{
	expression := &ast.InfixExpression {
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	
	return expression
}

func (p *Parser) parserExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(tok.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(tok.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(tok.RPAREN) {
		return nil
	}
	return exp
}


func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(tok.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tok.RPAREN) {
		return nil
	}
	if !p.expectPeek(tok.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(tok.ELSE) {
		p.nextToken()
		if !p.expectPeek(tok.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}


func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(tok.RBRACE) && !p.curTokenIs(tok.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}