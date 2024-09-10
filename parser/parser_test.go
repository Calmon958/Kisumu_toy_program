package parser

import (
	"fmt"
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
	checkParserErrors(t, p)

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
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lex.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. Got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression not *ast.Identifier. Got=%T", stmt.Expression)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. Got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntergerLiteralExpression(t *testing.T) {
	input := "9"

	l := lex.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. Got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression not *ast.IntegerLiteral. Got=%T", stmt.Expression)
	}
	if literal.Value != 9 {
		t.Errorf("literal.Value not %d. Got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "9" {
		t.Errorf("literal.TokenLiteral not %s. Got=%s", "9", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, ch := range prefixTests {
		l := lex.New(ch.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements doesn't contain %d statement. Got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. Got=%T", stmt.Expression)
		}
		if exp.Operator != ch.operator {
			t.Fatalf("exp.Operator is not '%s'. Got=%s", ch.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, ch.intValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	nb, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. Got=%T", il)
		return false
	}
	if nb.Value != value {
		t.Errorf("nb.Value not %d. Got=%d", value, nb.Value)
	}
	if nb.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("nb.TokenLiteral not %d. Got=%s", value, nb.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lex.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. Got=%T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftVal) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. Got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightVal) {
			return
		}
	}
}


func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},


		
	}

	for _, tt := range tests {
		l := lex.New(tt.input)
p := New(l)
program := p.ParseProgram()
checkParserErrors(t,p)
actual := program.String()
if actual != tt.expected {
	t.Errorf("expected=%q, Got=%q", tt.expected, actual)
}
	}
}