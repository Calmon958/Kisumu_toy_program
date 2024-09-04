package parser

import (
	"testing"
	lex "token/lexer"

	ast "token/ast"
)

func TestLetStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 918273;
	`
	l := lex.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t,p)
	
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't have 3 statements. Got %d", len(program.Statements))
	}
for _, stmt := range program.Statements {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt not *ast.returnStatement. got = %T", stmt)
		continue
	}
	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
	}
}
}



// checks for parser errors and prints them and kill the execution if it encounters an error
func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors{
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}