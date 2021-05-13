package fdalang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestParse(t *testing.T) {
	input := `a = 5 + 6
b = 3
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	vars := []string{"a", "b"}
	for i, stmt := range astProgram.Statements {
		assert.IsType(t, &AstStatementWithVoidedExpression{}, stmt, "%d statement", i)
		assignStmt, _ := stmt.(*AstStatementWithVoidedExpression)
		assert.IsType(t, &AstAssignment{}, assignStmt.Expr, "%d statement", i)
		assignExpr, _ := assignStmt.Expr.(*AstAssignment)
		assert.Equal(t, assignExpr.Left.Value, vars[i], "%d statement", i)
	}
}

func TestParseUnary(t *testing.T) {
	input := `a = -5
b = -a
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	vars := []string{"a", "b"}
	for i, stmt := range astProgram.Statements {
		assert.IsType(t, &AstStatementWithVoidedExpression{}, stmt, "%d statement", i)

		assignStmt, _ := stmt.(*AstStatementWithVoidedExpression)
		assignExpr, _ := assignStmt.Expr.(*AstAssignment)
		assert.Equal(t, assignExpr.Left.Value, vars[i], "%d statement", i)

		assert.IsType(t, &AstUnary{}, assignExpr.Value, "%d statement", i)
	}
}

func TestParseReal(t *testing.T) {
	input := `a = 5.6
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 1)
	assert.IsType(t, &AstStatementWithVoidedExpression{}, astProgram.Statements[0])
	assignStmt, _ := astProgram.Statements[0].(*AstStatementWithVoidedExpression)
	assert.IsType(t, &AstAssignment{}, assignStmt.Expr)
	assignExpr, _ := assignStmt.Expr.(*AstAssignment)
	assert.IsType(t, &AstNumFloat{}, assignExpr.Value)
}

func TestParseFunctionAndFunctionCall(t *testing.T) {
	input := `a = fn() int {
   return 2
}
c = a()
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	assert.IsType(t, &AstStatementWithVoidedExpression{}, astProgram.Statements[0])
	assignStmt, _ := astProgram.Statements[0].(*AstStatementWithVoidedExpression)
	assert.IsType(t, &AstAssignment{}, assignStmt.Expr)
	assignExpr, _ := assignStmt.Expr.(*AstAssignment)
	assert.IsType(t, &AstFunction{}, assignExpr.Value)

	function, _ := assignExpr.Value.(*AstFunction)
	require.Len(t, function.StatementsBlock.Statements, 1)
	assert.IsType(t, &AstReturn{}, function.StatementsBlock.Statements[0])

	returnStmt, _ := function.StatementsBlock.Statements[0].(*AstReturn)
	assert.IsType(t, &AstNumInt{}, returnStmt.ReturnValue)

	assert.IsType(t, &AstStatementWithVoidedExpression{}, astProgram.Statements[1])
	assignStmt2, _ := astProgram.Statements[1].(*AstStatementWithVoidedExpression)
	assert.IsType(t, &AstAssignment{}, assignStmt2.Expr)
	assignExpr2, _ := assignStmt2.Expr.(*AstAssignment)
	assert.IsType(t, &AstFunctionCall{}, assignExpr2.Value)
}

func TestParseFunctionAndFunctionCallWithArgs(t *testing.T) {
	input := `a = fn(int x, int y) int {
   return x + y
}
c = a(2, 5)
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	assert.IsType(t, &AstStatementWithVoidedExpression{}, astProgram.Statements[0])
	assignStmt, _ := astProgram.Statements[0].(*AstStatementWithVoidedExpression)
	assert.IsType(t, &AstAssignment{}, assignStmt.Expr)
	assignExpr, _ := assignStmt.Expr.(*AstAssignment)
	assert.IsType(t, &AstFunction{}, assignExpr.Value)

	function, _ := assignExpr.Value.(*AstFunction)
	require.Len(t, function.StatementsBlock.Statements, 1)
	assert.IsType(t, &AstReturn{}, function.StatementsBlock.Statements[0])
	assert.Len(t, function.Arguments, 2)
	assert.Equal(t, "int", function.Arguments[0].VarType)
	assert.Equal(t, "int", function.Arguments[1].VarType)
	assert.Equal(t, "x", function.Arguments[0].Var.Value)
	assert.Equal(t, "y", function.Arguments[1].Var.Value)

	returnStmt, _ := function.StatementsBlock.Statements[0].(*AstReturn)
	assert.IsType(t, &AstBinOperation{}, returnStmt.ReturnValue)

	binExpression, _ := returnStmt.ReturnValue.(*AstBinOperation)
	assert.IsType(t, &AstIdentifier{}, binExpression.Left)
	assert.IsType(t, &AstIdentifier{}, binExpression.Right)

	assert.IsType(t, &AstStatementWithVoidedExpression{}, astProgram.Statements[1])
	assignStmt2, _ := astProgram.Statements[1].(*AstStatementWithVoidedExpression)
	assert.IsType(t, &AstAssignment{}, assignStmt2.Expr)
	assignExpr2, _ := assignStmt2.Expr.(*AstAssignment)
	assert.IsType(t, &AstFunctionCall{}, assignExpr2.Value)
}

func TestParseIfStatement(t *testing.T) {
	input := `if 2 > 3 {
a = 4
}
b = 2
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	assert.IsType(t, &AstIfStatement{}, astProgram.Statements[0])

	ifStatement, _ := astProgram.Statements[0].(*AstIfStatement)
	assert.NotNil(t, ifStatement.Condition)
}

func TestParseIfStatementWithElseBranch(t *testing.T) {
	input := `if 2 > 3 {
a = 4
} else {
c = 3
}
b = 2
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)

	astProgram, err := p.Parse()
	require.Nil(t, err)

	require.Len(t, astProgram.Statements, 2)
	assert.IsType(t, &AstIfStatement{}, astProgram.Statements[0])
}

func TestArrayAsInvalidStatementNegative(t *testing.T) {
	input := `int[]{1, 2.1, 3}
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)
	_, err = p.Parse()
	require.NotNil(t, err)
}

func TestBinExprAsInvalidStatementNegative(t *testing.T) {
	input := `5 + 10
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)
	_, err = p.Parse()
	require.NotNil(t, err)
}

func TestInvalidExpr(t *testing.T) {
	input := `a = 5 + 10 )
`
	l := NewLexer(input)
	p, err := NewParser(l)
	require.Nil(t, err)
	_, err = p.Parse()
	require.NotNil(t, err)
}
