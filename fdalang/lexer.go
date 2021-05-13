package fdalang

import (
	"errors"
	"fmt"
	"unicode"
)

type Lexer struct {
	input        []rune
	currPosition int
	currChar     rune
	nextChar     rune
	line         int
	pos          int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}

	l.fetch(1, 1)
	return l
}

func (l *Lexer) fetch(line, pos int) {
	l.currChar = l.input[l.currPosition]
	l.nextChar = l.input[l.currPosition+1]
	l.line = line
	l.pos = pos
}

func ParseString(s string) ([]Token, error) {
	l := NewLexer(s)
	tokens := make([]Token, 0)
	for t, err := l.NextToken(); t.ID != TokenEOC; {
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (l *Lexer) GetCurrentPosition() int {
	return l.currPosition
}

func (l *Lexer) read() {
	l.currPosition += 1
	l.currChar = l.nextChar
	if l.currPosition+1 >= len(l.input) {
		l.nextChar = rune(0)
	} else {
		l.nextChar = l.input[l.currPosition+1]
	}

	if l.currPosition+1 < len(l.input) && l.input[l.currPosition-1] == '\n' {
		l.line += 1
		l.pos = 1
	} else {
		l.pos += 1
	}
}

func (l *Lexer) BackToToken(t Token) {
	l.currPosition = t.Pos
	l.fetch(t.Line, t.Col)
}

func (l *Lexer) NextToken() (Token, error) {
	var currToken Token
	l.skipWhitespace()

	currToken.Line = l.line
	currToken.Col = l.pos
	currToken.Pos = l.currPosition

	simpleTokens := []TokenID{
		TokenComma,
		TokenColon,
		TokenQuestion,
		TokenDot,
		TokenPlus,
		TokenMinus,
		TokenAsterisk,
		TokenLParen,
		TokenRParen,
		TokenLBrace,
		TokenRBrace,
		TokenLBracket,
		TokenRBracket,
		TokenLt,
		TokenGt,
	}
	for _, simpleToken := range simpleTokens {
		if string(l.currChar) == string(simpleToken) {
			currToken.ID = simpleToken
			currToken.Value = string(l.currChar)
			l.read()
			return currToken, nil
		}
	}

	switch l.currChar {
	case '\n':
		currToken.Value = ""
		currToken.ID = TokenEOL
	case '=':
		if l.nextChar == '=' {
			currToken.ID = TokenEq
			currToken.Value = string(TokenEq)
			l.read()
		} else {
			currToken.ID = TokenAssignment
			currToken.Value = string(TokenAssignment)
		}
	case '!':
		if l.nextChar == '=' {
			currToken.ID = TokenNotEq
			currToken.Value = string(TokenNotEq)
			l.read()
		} else {
			currToken.ID = TokenNot
			currToken.Value = string(TokenNot)
		}
	case '&':
		if l.nextChar != '&' {
			return currToken, l.error("Unexpected one `&`. Did you mean '&&'?")
		}
		currToken.ID = TokenAnd
		currToken.Value = string(TokenAnd)
		l.read()
	case '|':
		if l.nextChar != '|' {
			return currToken, l.error("Unexpected one `|`. Did you mean '||'?")
		}
		currToken.ID = TokenOr
		currToken.Value = string(TokenOr)
		l.read()
	case '/':
		if l.nextChar == '/' {
			l.consumeComment()
			return l.NextToken()
		} else {
			currToken.ID = TokenSlash
			currToken.Value = string(TokenSlash)
		}
	case 0:
		currToken.Value = ""
		currToken.ID = TokenEOC
	default:
		if isDigit(l.currChar) {
			value, isInt := l.readNumber()
			currToken.Value = value
			if isInt {
				currToken.ID = TokenNumInt
			} else {
				currToken.ID = TokenNumFloat
			}
		} else if unicode.IsLetter(l.currChar) {
			currToken.Value = l.readIdentifier()
			currToken.ID = LookupIdent(currToken.Value)
		} else {
			return currToken, l.error("Unexpected symbol: '%c'", l.currChar)
		}
	}
	l.read()
	return currToken, nil
}

func (l *Lexer) error(format string, args ...interface{}) error {
	errorMsg := fmt.Sprintf(format, args...)
	return errors.New(fmt.Sprintf("%s\nline:%d, pos %d", errorMsg, l.line, l.pos))
}

func (l *Lexer) GetCurrLineAndPos() (int, int) {
	return l.line, l.pos
}

func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' {
		l.read()
	}
}

func (l *Lexer) consumeComment() {
	for l.currChar != '\n' {
		l.read()
	}
}

func (l *Lexer) readNumber() (string, bool) {
	isInt := true
	result := string(l.currChar)
	for isDigit(l.nextChar) {
		result += string(l.nextChar)
		l.read()
	}
	if l.nextChar == '.' {
		isInt = false
		l.read()
		result += "."
		for isDigit(l.nextChar) {
			result += string(l.nextChar)
			l.read()
		}
	}

	return result, isInt
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	result := string(l.currChar)
	for unicode.IsLetter(l.nextChar) || isDigit(l.nextChar) {
		result += string(l.nextChar)
		l.read()
	}
	return result
}
