package fdalang

import (
	"fmt"
	"strings"
)

const (
	TokenSOL TokenID = "start of line"
	TokenEOL TokenID = "enf of line"
	TokenEOC TokenID = "end of code"

	TokenAssignment TokenID = "="
	TokenComma      TokenID = ","
	TokenDot        TokenID = "."
	TokenColon      TokenID = ":"
	TokenQuestion   TokenID = "?"

	// arithmetical operators
	TokenPlus     TokenID = "+"
	TokenMinus    TokenID = "-"
	TokenAsterisk TokenID = "*"
	TokenSlash    TokenID = "/"

	// logical operators
	TokenLt    TokenID = "<"
	TokenGt    TokenID = ">"
	TokenEq    TokenID = "=="
	TokenNotEq TokenID = "!="
	TokenNot   TokenID = "!"
	TokenAnd   TokenID = "&&"
	TokenOr    TokenID = "||"

	TokenNumInt   TokenID = "int"
	TokenNumFloat TokenID = "float"

	TokenLParen   TokenID = "("
	TokenRParen   TokenID = ")"
	TokenLBrace   TokenID = "{"
	TokenRBrace   TokenID = "}"
	TokenLBracket TokenID = "["
	TokenRBracket TokenID = "]"

	TokenIdent TokenID = "ident"

	// keywords
	TokenStruct   TokenID = "struct"
	TokenEnum     TokenID = "enum"
	TokenFunction TokenID = "fn"
	TokenReturn   TokenID = "return"
	TokenTrue     TokenID = "true"
	TokenFalse    TokenID = "false"
	TokenIf       TokenID = "if"
	TokenElse     TokenID = "else"
	TokenSwitch   TokenID = "switch"
	TokenCase     TokenID = "case"
	TokenDefault  TokenID = "default"

	// type hints
	TokenType TokenID = "type"

	TokenInvalid TokenID = "invalid"
)

type TokenID string

type Token struct {
	ID    TokenID
	Value string
	Line  int
	Col   int
	Pos   int
}

var strToKeywordMap = map[string]TokenID{
	"fn":      TokenFunction,
	"return":  TokenReturn,
	"void":    TokenType,
	"int":     TokenType,
	"float":   TokenType,
	"true":    TokenTrue,
	"false":   TokenFalse,
	"if":      TokenIf,
	"else":    TokenElse,
	"struct":  TokenStruct,
	"enum":    TokenEnum,
	"switch":  TokenSwitch,
	"case":    TokenCase,
	"default": TokenDefault,
}

func TokensKeywords() map[TokenID]bool {
	return map[TokenID]bool{
		TokenFunction: true,
		TokenReturn: true,
		TokenType: true,
		TokenTrue: true,
		TokenFalse: true,
		TokenIf: true,
		TokenElse: true,
		TokenStruct: true,
		TokenEnum: true,
		TokenSwitch: true,
		TokenCase: true,
		TokenDefault: true,
	}
}

func TokensConstants() map[TokenID]bool {
	return map[TokenID]bool{
		TokenNumInt: true,
		TokenNumFloat: true,
		TokenTrue: true,
		TokenFalse: true,
	}
}

func keywordOrIdent(ident string) TokenID {
	if k, ok := strToKeywordMap[ident]; ok {
		return k
	}

	return TokenIdent
}

func TokenIDs(tokens TokenID) []TokenID {
	return []TokenID{tokens}
}

func TokensString(tokens []TokenID) string {
	var s []string
	for _, t := range tokens {
		s = append(s, fmt.Sprintf("'%s'", t))
	}
	return strings.Join(s, ", ")
}
