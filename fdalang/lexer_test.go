package fdalang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type expectedTestToken struct {
	expectedType  TokenID
	expectedValue string
}

func TestNextTokenGeneric(t *testing.T) {
	input := `a = (5 + 6)
b = 3 > 2 < 1
// comment here
v = b == 1
z = b != 1
x = !y == true
c = fn(int a, int b) int {
   return 3 + a
}`

	tests := []expectedTestToken{
		{TokenIdent, "a"},
		{TokenAssignment, "="},
		{TokenLParen, "("},
		{TokenNumInt, "5"},
		{TokenPlus, "+"},
		{TokenNumInt, "6"},
		{TokenRParen, ")"},
		{TokenEOL, ""},
		{TokenIdent, "b"},
		{TokenAssignment, "="},
		{TokenNumInt, "3"},
		{TokenGt, ">"},
		{TokenNumInt, "2"},
		{TokenLt, "<"},
		{TokenNumInt, "1"},
		{TokenEOL, ""},
		{TokenEOL, ""},
		{TokenIdent, "v"},
		{TokenAssignment, "="},
		{TokenIdent, "b"},
		{TokenEq, "=="},
		{TokenNumInt, "1"},
		{TokenEOL, ""},
		{TokenIdent, "z"},
		{TokenAssignment, "="},
		{TokenIdent, "b"},
		{TokenNotEq, "!="},
		{TokenNumInt, "1"},
		{TokenEOL, ""},
		{TokenIdent, "x"},
		{TokenAssignment, "="},
		{TokenNot, "!"},
		{TokenIdent, "y"},
		{TokenEq, "=="},
		{TokenTrue, "true"},
		{TokenEOL, ""},
		{TokenIdent, "c"},
		{TokenAssignment, "="},
		{TokenFunction, "fn"},
		{TokenLParen, "("},
		{TokenType, "int"},
		{TokenIdent, "a"},
		{TokenComma, ","},
		{TokenType, "int"},
		{TokenIdent, "b"},
		{TokenRParen, ")"},
		{TokenType, "int"},
		{TokenLBrace, "{"},
		{TokenEOL, ""},
		{TokenReturn, "return"},
		{TokenNumInt, "3"},
		{TokenPlus, "+"},
		{TokenIdent, "a"},
		{TokenEOL, ""},
		{TokenRBrace, "}"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestReal(t *testing.T) {
	input := `a = 5.6`

	tests := []expectedTestToken{
		{TokenIdent, "a"},
		{TokenAssignment, "="},
		{TokenNumFloat, "5.6"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestArray(t *testing.T) {
	input := `arr = int[]{1, 2}
o = arr[0]`
	tests := []expectedTestToken{
		{TokenIdent, "arr"},
		{TokenAssignment, "="},
		{TokenType, "int"},
		{TokenLBracket, "["},
		{TokenRBracket, "]"},
		{TokenLBrace, "{"},
		{TokenNumInt, "1"},
		{TokenComma, ","},
		{TokenNumInt, "2"},
		{TokenRBrace, "}"},
		{TokenEOL, ""},
		{TokenIdent, "o"},
		{TokenAssignment, "="},
		{TokenIdent, "arr"},
		{TokenLBracket, "["},
		{TokenNumInt, "0"},
		{TokenRBracket, "]"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestRealShort(t *testing.T) {
	input := `a = 5.`

	tests := []expectedTestToken{
		{TokenIdent, "a"},
		{TokenAssignment, "="},
		{TokenNumFloat, "5."},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestLogicalAndOr(t *testing.T) {
	input := `a = true && false || false`

	tests := []expectedTestToken{
		{TokenIdent, "a"},
		{TokenAssignment, "="},
		{TokenTrue, "true"},
		{TokenAnd, "&&"},
		{TokenFalse, "false"},
		{TokenOr, "||"},
		{TokenFalse, "false"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestLexerStruct(t *testing.T) {
	input := `struct point {
   float x
   float y
}
p = point{x = 1., y = 2.}
px = p.x`

	tests := []expectedTestToken{
		{TokenStruct, "struct"},
		{TokenIdent, "point"},
		{TokenLBrace, "{"},
		{TokenEOL, ""},
		{TokenType, "float"},
		{TokenIdent, "x"},
		{TokenEOL, ""},
		{TokenType, "float"},
		{TokenIdent, "y"},
		{TokenEOL, ""},
		{TokenRBrace, "}"},
		{TokenEOL, ""},
		{TokenIdent, "p"},
		{TokenAssignment, "="},
		{TokenIdent, "point"},
		{TokenLBrace, "{"},
		{TokenIdent, "x"},
		{TokenAssignment, "="},
		{TokenNumFloat, "1."},
		{TokenComma, ","},
		{TokenIdent, "y"},
		{TokenAssignment, "="},
		{TokenNumFloat, "2."},
		{TokenRBrace, "}"},
		{TokenEOL, ""},
		{TokenIdent, "px"},
		{TokenAssignment, "="},
		{TokenIdent, "p"},
		{TokenDot, "."},
		{TokenIdent, "x"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestSwitchCase(t *testing.T) {
	input := `switch {
case a == 1:
   r = 1
default:
   r = 2
}`

	tests := []expectedTestToken{
		{TokenSwitch, "switch"},
		{TokenLBrace, "{"},
		{TokenEOL, ""},
		{TokenCase, "case"},
		{TokenIdent, "a"},
		{TokenEq, "=="},
		{TokenNumInt, "1"},
		{TokenColon, ":"},
		{TokenEOL, ""},
		{TokenIdent, "r"},
		{TokenAssignment, "="},
		{TokenNumInt, "1"},
		{TokenEOL, ""},
		{TokenDefault, "default"},
		{TokenColon, ":"},
		{TokenEOL, ""},
		{TokenIdent, "r"},
		{TokenAssignment, "="},
		{TokenNumInt, "2"},
		{TokenEOL, ""},
		{TokenRBrace, "}"},
		{TokenEOC, ""},
	}

	testLexerInput(input, tests, t)
}

func TestGetCurrLineAndPos(t *testing.T) {
	input := `a = 5 + 6
asd`
	l := NewLexer(input)
	assert.Equal(t, 1, l.line, "Line should be 1 on start")
	assert.Equal(t, 1, l.pos, "Col should be 1 on start")

	l.read()
	assert.Equal(t, 1, l.line, "Line should be 1")
	assert.Equal(t, 2, l.pos, "Col should be 2")

	for i := 0; i <= 8; i++ {
		l.read()
	}
	assert.Equal(t, 2, l.line, "Line should be 2")
	assert.Equal(t, 1, l.pos, "Col should be 1")
}

func TestLineAndPosForTokens(t *testing.T) {
	input := `a = fn() {
   b = 3
}`
	l := NewLexer(input)
	_, _ = l.NextToken()
	_, _ = l.NextToken()
	tok, _ := l.NextToken()
	assert.Equal(t, 1, tok.Line)
	assert.Equal(t, 5, tok.Col)

	_, _ = l.NextToken()
	_, _ = l.NextToken()
	_, _ = l.NextToken()
	tok, _ = l.NextToken()
	assert.Equal(t, 1, tok.Line)
	assert.Equal(t, 11, tok.Col)

	tok, _ = l.NextToken()
	assert.Equal(t, 2, tok.Line)
	assert.Equal(t, 4, tok.Col)
}

func testLexerInput(input string, tests []expectedTestToken, t *testing.T) {
	l := NewLexer(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err, "[%d] token lexer error", i)
		require.Equal(t, tt.expectedType, tok.Type, "[%d] token type wrong", i)
		require.Equal(t, tt.expectedValue, tok.Value, "[%d] token value wrong", i)
	}
}
