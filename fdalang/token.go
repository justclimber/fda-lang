package fdalang

import (
	"fmt"
	"strings"
)

const (
	TokenEOL = "TokenEOL"
	TokenEOF = "TokenEOF"

	TokenAssignment = "="
	TokenComma      = ","
	TokenDot        = "."
	TokenColon      = ":"
	TokenQuestion   = "?"

	// arithmetical operators
	TokenPlus     = "+"
	TokenMinus    = "-"
	TokenAsterisk = "*"
	TokenSlash    = "/"

	// logical operators
	TokenLt    = "<"
	TokenGt    = ">"
	TokenEq    = "=="
	TokenNotEq = "!="
	TokenNot   = "!"
	TokenAnd   = "&&"
	TokenOr    = "||"

	TokenNumInt   = "int_num"
	TokenNumFloat = "float_num"

	TokenLParen   = "("
	TokenRParen   = ")"
	TokenLBrace   = "{"
	TokenRBrace   = "}"
	TokenLBracket = "["
	TokenRBracket = "]"

	TokenIdent = "ident"

	// keywords
	TokenStruct   = "struct"
	TokenEnum     = "enum"
	TokenFunction = "fn"
	TokenReturn   = "return"
	TokenTrue     = "true"
	TokenFalse    = "false"
	TokenIf       = "if"
	TokenElse     = "else"
	TokenSwitch   = "switch"
	TokenCase     = "case"
	TokenDefault  = "default"

	// type hints
	TokenType = "type"
)

type TokenID string

type Token struct {
	Type  TokenID
	Value string
	Line  int
	Col   int
	Pos   int
}

var keywords = map[string]TokenID{
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

func LookupIdent(ident string) TokenID {
	if keywordToken, ok := keywords[ident]; ok {
		return keywordToken
	}

	return TokenIdent
}

func TokenIDs(tokens TokenID) []TokenID {
	return []TokenID{tokens}
}

func GetTokensString(tokens []TokenID) string {
	var s []string
	for _, t := range tokens {
		s = append(s, fmt.Sprintf("'%s'", t))
	}
	return strings.Join(s, ", ")
}
