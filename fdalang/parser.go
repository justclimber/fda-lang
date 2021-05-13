package fdalang

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	PriorityLowest
	PriorityAssignment // =
	PriorityOr         // ||
	PriorityAnd        // &&
	PriorityEquals     // ==
	PriorityComparison // > or <
	PrioritySum        // +
	PriorityProduct    // *
	PriorityPrefix     // -X or !X
	PriorityCall       // myFunction(X)
	PriorityIndex      // array[index]
)

var precedences = map[TokenID]int{
	TokenEq:         PriorityEquals,
	TokenNotEq:      PriorityEquals,
	TokenLt:         PriorityComparison,
	TokenGt:         PriorityComparison,
	TokenAssignment: PriorityAssignment,
	TokenAnd:        PriorityAnd,
	TokenOr:         PriorityOr,
	TokenPlus:       PrioritySum,
	TokenMinus:      PrioritySum,
	TokenSlash:      PriorityProduct,
	TokenAsterisk:   PriorityProduct,
	TokenLParen:     PriorityCall,
	TokenLBracket:   PriorityIndex,
	TokenLBrace:     PriorityIndex,
	TokenDot:        PriorityIndex,
	TokenColon:      PriorityIndex,
}

type (
	unaryExprFunction func([]TokenID) (IExpression, error)
	binExprFunctions  func(IExpression, []TokenID) (IExpression, error)
)

type Parser struct {
	l *Lexer

	currToken Token
	nextToken Token
	prevToken Token

	unaryExprFunctions map[TokenID]unaryExprFunction
	binExprFunctions   map[TokenID]binExprFunctions
}

func NewParser(l *Lexer) (*Parser, error) {
	p := &Parser{l: l}

	var err error
	p.currToken, err = p.l.NextToken()
	if err != nil {
		return nil, err
	}

	p.nextToken, err = p.l.NextToken()
	if err != nil {
		return nil, err
	}

	p.unaryExprFunctions = make(map[TokenID]unaryExprFunction)
	p.registerUnaryExprFunction(TokenMinus, p.parseUnaryExpression)
	p.registerUnaryExprFunction(TokenNot, p.parseUnaryExpression)
	p.registerUnaryExprFunction(TokenNumInt, p.parseInteger)
	p.registerUnaryExprFunction(TokenNumFloat, p.parseReal)
	p.registerUnaryExprFunction(TokenTrue, p.parseBoolean)
	p.registerUnaryExprFunction(TokenFalse, p.parseBoolean)
	p.registerUnaryExprFunction(TokenIdent, p.parseIdentifierAsExpression)
	p.registerUnaryExprFunction(TokenLParen, p.parseGroupedExpression)
	p.registerUnaryExprFunction(TokenFunction, p.parseFunction)
	p.registerUnaryExprFunction(TokenLBracket, p.parseArray)
	p.registerUnaryExprFunction(TokenQuestion, p.parseEmptierExpression)

	p.binExprFunctions = make(map[TokenID]binExprFunctions)
	p.registerBinExprFunction(TokenPlus, p.parseBinExpression)
	p.registerBinExprFunction(TokenMinus, p.parseBinExpression)
	p.registerBinExprFunction(TokenSlash, p.parseBinExpression)
	p.registerBinExprFunction(TokenLt, p.parseBinExpression)
	p.registerBinExprFunction(TokenGt, p.parseBinExpression)
	p.registerBinExprFunction(TokenEq, p.parseBinExpression)
	p.registerBinExprFunction(TokenAnd, p.parseBinExpression)
	p.registerBinExprFunction(TokenOr, p.parseBinExpression)
	p.registerBinExprFunction(TokenNotEq, p.parseBinExpression)
	p.registerBinExprFunction(TokenAsterisk, p.parseBinExpression)
	p.registerBinExprFunction(TokenLParen, p.parseFunctionCall)
	p.registerBinExprFunction(TokenLBracket, p.parseArrayIndexCall)
	p.registerBinExprFunction(TokenLBrace, p.parseStructExpression)
	p.registerBinExprFunction(TokenDot, p.parseStructFieldCall)
	p.registerBinExprFunction(TokenColon, p.parseEnumExpression)

	return p, nil
}

func (p *Parser) read() error {
	var err error
	p.prevToken = p.currToken
	p.currToken = p.nextToken
	p.nextToken, err = p.l.NextToken()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) readWithEolOpt() error {
	if err := p.read(); err != nil {
		return err
	}
	if p.currToken.Type == TokenEOL {
		if err := p.read(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) back() {
	p.l.BackToToken(p.prevToken)
	p.nextToken = p.currToken
	p.currToken = p.prevToken
	_ = p.read()
	_ = p.read()
}

func (p *Parser) Parse() (*StatementsBlock, error) {
	program := &StatementsBlock{}

	statements, err := p.parseBlockOfStatements(TokenIDs(TokenEOF))
	program.Statements = statements

	return program, err
}

func (p *Parser) parseBlockOfStatements(terminatedTokens []TokenID) ([]IStatement, error) {
	var statements []IStatement

	for !p.currTokenIn(terminatedTokens) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
		if err = p.read(); err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (p *Parser) parseStatement() (IStatement, error) {
	switch p.currToken.Type {
	case TokenIdent:
		if p.nextToken.Type == TokenLParen {
			function := &Identifier{
				Token: p.currToken,
				Value: p.currToken.Value,
			}
			err := p.read()
			if err != nil {
				return nil, err
			}
			return p.parseFunctionCall(function, TokenIDs(TokenEOL))
		} else if p.nextToken.Type == TokenDot {
			return p.parseStructFieldAssignment(TokenIDs(TokenEOL))
		} else {
			return p.parseAssignment(TokenIDs(TokenEOL))
		}
	case TokenReturn:
		return p.parseReturn()
	case TokenIf:
		return p.parseIfStatement()
	case TokenStruct:
		return p.parseStructDefinition()
	case TokenEnum:
		return p.parseEnumDefinition()
	case TokenSwitch:
		return p.parseSwitchStatement()
	case TokenEOL:
		return nil, nil
	default:
		return nil, p.parseError("Unexpected token for start of statement: %s\n", p.currToken.Type)
	}
}

func (p *Parser) parseStructFieldAssignment(terminatedTokens []TokenID) (*StructFieldAssignment, error) {
	assignStmt := &StructFieldAssignment{Token: p.currToken}

	left, err := p.parseIdentifier(terminatedTokens)
	if err != nil {
		return nil, err
	}

	var leftWithFieldCall IExpression
	leftWithFieldCall = left

	// nested structs can be here
	for p.nextToken.Type == TokenDot {
		if err = p.read(); err != nil {
			return nil, err
		}

		if _, err = p.getExpectedToken(TokenDot); err != nil {
			return nil, err
		}

		leftWithFieldCall, err = p.parseStructFieldCall(leftWithFieldCall, terminatedTokens)
		if err != nil {
			return nil, err
		}
	}
	assignStmt.Left = leftWithFieldCall.(*StructFieldCall)

	if err = p.read(); err != nil {
		return nil, err
	}

	if _, err = p.getExpectedToken(TokenAssignment); err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}
	assignStmt.Value, err = p.parseExpression(PriorityLowest, terminatedTokens)
	if err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}

	if _, err = p.getExpectedTokens(terminatedTokens); err != nil {
		return nil, err
	}

	return assignStmt, nil
}

func (p *Parser) parseAssignment(terminatedTokens []TokenID) (*Assignment, error) {
	assignStmt := &Assignment{Token: p.currToken}
	identStmt, err := p.parseIdentifier(terminatedTokens)
	if err != nil {
		return nil, err
	}
	assignStmt.Left = identStmt

	if err = p.read(); err != nil {
		return nil, err
	}
	if _, err = p.getExpectedToken(TokenAssignment); err != nil {
		return nil, err
	}
	if err = p.read(); err != nil {
		return nil, err
	}
	assignStmt.Value, err = p.parseExpression(PriorityLowest, terminatedTokens)
	if err != nil {
		return nil, err
	}
	if err = p.read(); err != nil {
		return nil, err
	}

	if _, err = p.getExpectedTokens(terminatedTokens); err != nil {
		return nil, err
	}

	return assignStmt, nil
}

func (p *Parser) parseReturn() (*Return, error) {
	stmt := &Return{Token: p.currToken}
	var err error
	if err = p.read(); err != nil {
		return nil, err
	}

	stmt.ReturnValue, err = p.parseExpression(PriorityLowest, TokenIDs(TokenEOL))

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int, terminatedTokens []TokenID) (IExpression, error) {
	unaryFunction := p.unaryExprFunctions[p.currToken.Type]
	if unaryFunction == nil {
		err := p.parseError("no Unary parse function for %s found", p.currToken.Type)
		return nil, err
	}

	leftExpr, err := unaryFunction(terminatedTokens)
	if err != nil {
		return nil, err
	}

	return p.parseRightPartOfExpression(leftExpr, precedence, terminatedTokens)
}

func (p *Parser) parseRightPartOfExpression(
	leftExpr IExpression,
	precedence int,
	terminatedTokens []TokenID,
) (IExpression, error) {
	var err error
	nextPrecedence := p.nextPrecedence()
	for !p.nextTokenIn(terminatedTokens) && precedence < nextPrecedence {
		binExprFunction := p.binExprFunctions[p.nextToken.Type]
		if binExprFunction == nil {
			err := p.parseError("Unexpected next token for binary expression '%s'", p.nextToken.Type)
			return nil, err
		}

		if err = p.read(); err != nil {
			return nil, err
		}
		leftExpr, err = binExprFunction(leftExpr, terminatedTokens)

		if err != nil {
			return nil, err
		}
	}
	return leftExpr, nil
}

func (p *Parser) parseIdentifierAsExpression(terminatedTokens []TokenID) (IExpression, error) {
	_, err := p.getExpectedToken(TokenIdent)
	if err != nil {
		return nil, err
	}
	return &Identifier{
		Token: p.currToken,
		Value: p.currToken.Value,
	}, nil
}

func (p *Parser) parseIdentifier(terminatedTokens []TokenID) (*Identifier, error) {
	expr, err := p.parseIdentifierAsExpression(terminatedTokens)
	if err != nil {
		return nil, err
	}
	ident, _ := expr.(*Identifier)
	return ident, nil
}

func (p *Parser) parseInteger(terminatedTokens []TokenID) (IExpression, error) {
	node := &NumInt{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Value, 0, 64)
	if err != nil {
		err := p.parseError("could not parse %q as integer", p.currToken.Value)
		return nil, err
	}

	node.Value = value
	return node, nil
}

func (p *Parser) parseUnaryExpression(terminatedTokens []TokenID) (IExpression, error) {
	node := &UnaryExpression{
		Token:    p.currToken,
		Operator: p.currToken.Value,
	}
	if err := p.read(); err != nil {
		return nil, err
	}
	expressionResult, err := p.parseExpression(PriorityPrefix, terminatedTokens)
	if err != nil {
		return nil, err
	}
	node.Right = expressionResult

	return node, err
}

func (p *Parser) parseEmptierExpression(terminatedTokens []TokenID) (IExpression, error) {
	node := &EmptierExpression{Token: p.currToken, IsArray: false}
	if err := p.read(); err != nil {
		return nil, err
	}

	_, err := p.getExpectedTokens([]TokenID{TokenLBracket, TokenType, TokenIdent})
	if err != nil {
		return nil, err
	}
	if p.currToken.Type == TokenLBracket {
		if err := p.requireToken(TokenRBracket); err != nil {
			return nil, err
		}
		node.IsArray = true
		if err := p.read(); err != nil {
			return nil, err
		}
	}

	node.Type = p.currToken.Value
	return node, nil
}

func (p *Parser) parseReal(terminatedTokens []TokenID) (IExpression, error) {
	node := &NumFloat{Token: p.currToken}

	value, err := strconv.ParseFloat(p.currToken.Value, 64)
	if err != nil {
		err := p.parseError("could not parse %q as float", p.currToken.Value)
		return nil, err
	}

	node.Value = value
	return node, nil
}

func (p *Parser) parseBinExpression(left IExpression, terminatedTokens []TokenID) (IExpression, error) {
	expression := &BinExpression{
		Token:    p.currToken,
		Operator: p.currToken.Value,
		Left:     left,
	}
	var err error
	precedence := p.curPrecedence()
	if err = p.read(); err != nil {
		return nil, err
	}

	expression.Right, err = p.parseExpression(precedence, terminatedTokens)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseSwitchStatement() (IExpression, error) {
	stmt := &Switch{Token: p.currToken}

	var err error
	if p.nextToken.Type != TokenLBrace {
		if err = p.read(); err != nil {
			return nil, err
		}

		expr, err := p.parseExpression(PriorityLowest, TokenIDs(TokenLBrace))
		if err != nil {
			return nil, err
		}
		stmt.SwitchExpression = expr
	}

	if err = p.requireTokenSequence([]TokenID{TokenLBrace, TokenEOL}); err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}

	cases := make([]*Case, 0)
	for p.currToken.Type == TokenCase {
		caseBlock := &Case{Token: Token{}}

		if stmt.SwitchExpression != nil {
			caseBlock.Condition, err = p.parseRightPartOfExpression(
				stmt.SwitchExpression,
				PriorityLowest,
				TokenIDs(TokenEOL),
			)
		} else {
			if err = p.read(); err != nil {
				return nil, err
			}
			caseBlock.Condition, err = p.parseExpression(PriorityLowest, TokenIDs(TokenEOL))
		}
		if err != nil {
			return nil, err
		}

		if err = p.requireToken(TokenEOL); err != nil {
			return nil, err
		}

		statements, err := p.parseBlockOfStatements([]TokenID{TokenCase, TokenDefault, TokenRBrace})
		if err != nil {
			return nil, err
		}
		caseBlock.PositiveBranch = &StatementsBlock{Statements: statements}
		cases = append(cases, caseBlock)
	}
	stmt.Cases = cases

	if p.currToken.Type == TokenDefault {
		if err = p.requireToken(TokenEOL); err != nil {
			return nil, err
		}
		statements, err := p.parseBlockOfStatements(TokenIDs(TokenRBrace))
		if err != nil {
			return nil, err
		}
		stmt.DefaultBranch = &StatementsBlock{Statements: statements}
	}

	return stmt, nil
}

func (p *Parser) parseIfStatement() (IExpression, error) {
	stmt := &IfStatement{Token: p.currToken}

	var err error

	if err = p.read(); err != nil {
		return nil, err
	}

	stmt.Condition, err = p.parseExpression(PriorityLowest, TokenIDs(TokenLBrace))
	if err != nil {
		return nil, err
	}

	if err := p.requireTokenSequence([]TokenID{TokenLBrace, TokenEOL}); err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}

	statements, err := p.parseBlockOfStatements(TokenIDs(TokenRBrace))
	if err != nil {
		return nil, err
	}

	stmt.PositiveBranch = &StatementsBlock{Statements: statements}

	if err = p.read(); err != nil {
		return nil, err
	}

	if p.currToken.Type != TokenElse {
		return stmt, nil
	}

	if err := p.requireTokenSequence([]TokenID{TokenLBrace, TokenEOL}); err != nil {
		return nil, err
	}

	statements, err = p.parseBlockOfStatements(TokenIDs(TokenRBrace))
	stmt.ElseBranch = &StatementsBlock{Statements: statements}

	return stmt, err
}

func (p *Parser) parseStructDefinition() (IExpression, error) {
	node := &StructDefinition{Token: p.currToken}

	if err := p.read(); err != nil {
		return nil, err
	}
	name, err := p.getExpectedToken(TokenIdent)
	if err != nil {
		return nil, err
	}
	node.Name = name.Value

	if err := p.requireTokenSequence([]TokenID{TokenLBrace, TokenEOL}); err != nil {
		return nil, err
	}

	if err := p.read(); err != nil {
		return nil, err
	}

	fields, err := p.parseVarAndTypes(TokenRBrace, TokenEOL)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, p.parseError("Struct should contain at least 1 field")
	}

	fieldsMap := make(map[string]*VarAndType)
	for _, field := range fields {
		fieldsMap[field.Var.Value] = field
	}

	node.Fields = fieldsMap

	return node, nil
}

func (p *Parser) parseFunction(terminatedTokens []TokenID) (IExpression, error) {
	function := &Function{Token: p.currToken}

	err := p.read()
	if err != nil {
		return nil, err
	}
	_, err = p.getExpectedToken(TokenLParen)
	if err != nil {
		return nil, err
	}

	err = p.read()
	if err != nil {
		return nil, err
	}
	function.Arguments, err = p.parseVarAndTypes(TokenRParen, TokenComma)
	if err != nil {
		return nil, err
	}

	_, err = p.getExpectedToken(TokenRParen)
	if err != nil {
		return nil, err
	}

	err = p.read()
	if err != nil {
		return nil, err
	}
	typeToken, err := p.getExpectedTokens([]TokenID{TokenType, TokenIdent})
	if err != nil {
		return nil, err
	}
	function.ReturnType = typeToken.Value

	if err := p.requireTokenSequence([]TokenID{TokenLBrace, TokenEOL}); err != nil {
		return nil, err
	}

	err = p.read()
	if err != nil {
		return nil, err
	}
	statements, err := p.parseBlockOfStatements(TokenIDs(TokenRBrace))
	function.StatementsBlock = &StatementsBlock{Statements: statements}

	return function, err
}

func (p *Parser) parseVarAndTypes(endToken TokenID, delimiterToken TokenID) ([]*VarAndType, error) {
	var err error
	vars := make([]*VarAndType, 0)

	for p.currTokenIn([]TokenID{TokenLBracket, TokenType, TokenIdent}) {
		argument := &VarAndType{Token: p.currToken}
		arrayTypePrefix := ""
		if p.currToken.Type == TokenLBracket {
			if err := p.requireToken(TokenRBracket); err != nil {
				return nil, err
			}
			arrayTypePrefix = "[]"
			if err = p.read(); err != nil {
				return nil, err
			}
		}
		argument.VarType = arrayTypePrefix + p.currToken.Value

		if err = p.read(); err != nil {
			return nil, err
		}

		argument.Var, err = p.parseIdentifier(TokenIDs(delimiterToken))
		if err != nil {
			return nil, err
		}

		vars = append(vars, argument)

		if p.nextToken.Type != endToken {
			err = p.read()
			if err != nil {
				return nil, err
			}
			_, err := p.getExpectedToken(delimiterToken)
			if err != nil {
				return nil, err
			}
		}
		if err = p.readWithEolOpt(); err != nil {
			return nil, err
		}
	}

	return vars, nil
}

func (p *Parser) parseFunctionCall(function IExpression, terminatedTokens []TokenID) (IExpression, error) {
	functionCall := &FunctionCall{
		Token:    p.currToken,
		Function: function,
	}
	var err error
	if err = p.read(); err != nil {
		return nil, err
	}

	functionCall.Arguments, err = p.parseExpressions(TokenIDs(TokenRParen))

	return functionCall, err
}

func (p *Parser) parseExpressions(closeTokens []TokenID) ([]IExpression, error) {
	var expressions []IExpression

	for !p.currTokenIn(closeTokens) {
		expression, err := p.parseExpression(PriorityLowest, append(closeTokens, TokenComma))
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
		if err = p.read(); err != nil {
			return nil, err
		}
		if p.currToken.Type == TokenComma {
			if err = p.readWithEolOpt(); err != nil {
				return nil, err
			}
		}

	}

	return expressions, nil
}

func (p *Parser) parseGroupedExpression(terminatedTokens []TokenID) (IExpression, error) {
	err := p.read()
	if err != nil {
		return nil, err
	}

	expression, err := p.parseExpression(PriorityLowest, TokenIDs(TokenRParen))
	if err != nil {
		return nil, err
	}
	if err = p.read(); err != nil {
		return nil, err
	}
	_, err = p.getExpectedToken(TokenRParen)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseBoolean(terminatedTokens []TokenID) (IExpression, error) {
	return &Boolean{
		Token: p.currToken,
		Value: p.currToken.Type == TokenTrue,
	}, nil
}

func (p *Parser) parseArray(terminatedTokens []TokenID) (IExpression, error) {
	node := &Array{Token: p.currToken}

	var err error
	if err = p.requireToken(TokenRBracket); err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}

	arrayTypeToken, err := p.getExpectedTokens([]TokenID{TokenIdent, TokenType})
	if err != nil {
		return nil, err
	}

	node.ElementsType = arrayTypeToken.Value

	if err = p.read(); err != nil {
		return nil, err
	}

	var elementExpressions []IExpression
	if p.currToken.Type == TokenLBrace {
		if err = p.read(); err != nil {
			return nil, err
		}
		elementExpressions, err = p.parseExpressions([]TokenID{TokenComma, TokenRBrace})
		if err != nil {
			return nil, err
		}
	}

	node.Elements = elementExpressions

	return node, nil
}

func (p *Parser) parseArrayIndexCall(array IExpression, terminatedTokens []TokenID) (IExpression, error) {
	node := &ArrayIndexCall{
		Token: p.currToken,
		Left:  array,
	}

	var err error
	if err = p.read(); err != nil {
		return nil, err
	}

	index, err := p.parseExpression(PriorityIndex, []TokenID{TokenRBracket})
	if err != nil {
		return nil, err
	}

	if err = p.read(); err != nil {
		return nil, err
	}
	node.Index = index

	return node, nil
}

func (p *Parser) parseStructExpression(
	expr IExpression,
	terminatedTokens []TokenID,
) (IExpression, error) {
	ident, ok := expr.(*Identifier)
	if !ok {
		return nil, p.parseError("Struct operator should only on identifiers, but '%T'", expr)
	}
	node := &Struct{
		Token: p.currToken,
		Ident: ident,
	}
	if err := p.read(); err != nil {
		return nil, err
	}

	fields := make([]*Assignment, 0)
	for p.currToken.Type == TokenIdent {
		field, err := p.parseAssignment([]TokenID{TokenComma, TokenRBrace})
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
		if p.currToken.Type == TokenComma {
			if err = p.readWithEolOpt(); err != nil {
				return nil, err
			}
		}
	}
	node.Fields = fields

	return node, nil
}

func (p *Parser) parseStructFieldCall(expr IExpression, terminatedTokens []TokenID) (IExpression, error) {
	node := &StructFieldCall{
		Token:      p.currToken,
		StructExpr: expr,
	}
	if err := p.read(); err != nil {
		return nil, err
	}
	field, err := p.parseIdentifier(terminatedTokens)
	if err != nil {
		return nil, err
	}

	node.Field = field

	return node, nil
}

func (p *Parser) parseEnumDefinition() (IExpression, error) {
	node := &EnumDefinition{Token: p.currToken}

	if err := p.read(); err != nil {
		return nil, err
	}
	name, err := p.getExpectedToken(TokenIdent)
	if err != nil {
		return nil, err
	}
	node.Name = name.Value

	if err := p.requireToken(TokenLBrace); err != nil {
		return nil, err
	}

	if err := p.readWithEolOpt(); err != nil {
		return nil, err
	}

	node.Elements = make([]string, 0)
	for p.currToken.Type != TokenRBrace {
		el, err := p.getExpectedToken(TokenIdent)
		if err != nil {
			return nil, err
		}
		node.Elements = append(node.Elements, el.Value)
		if err := p.read(); err != nil {
			return nil, err
		}

		if p.currToken.Type == TokenComma {
			if err := p.readWithEolOpt(); err != nil {
				return nil, err
			}
		}
	}

	if err := p.read(); err != nil {
		return nil, err
	}

	return node, nil
}

func (p *Parser) parseEnumExpression(expr IExpression, terminatedTokens []TokenID) (IExpression, error) {
	node := &EnumElementCall{
		Token:    p.currToken,
		EnumExpr: expr,
	}
	if err := p.read(); err != nil {
		return nil, err
	}
	el, err := p.parseIdentifier(terminatedTokens)
	if err != nil {
		return nil, err
	}

	node.Element = el
	return node, nil
}

func (p *Parser) nextPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}

	return PriorityLowest
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return PriorityLowest
}

func (p *Parser) getExpectedToken(TokenID TokenID) (Token, error) {
	if p.currToken.Type != TokenID {
		err := p.parseError("expected token to be '%s', got '%s' instead",
			TokenID, p.currToken.Type)
		return Token{}, err
	}
	return p.currToken, nil
}

func (p *Parser) getExpectedTokens(tokenTypes []TokenID) (Token, error) {
	if len(tokenTypes) == 1 {
		return p.getExpectedToken(tokenTypes[0])
	}
	for _, tok := range tokenTypes {
		if p.currToken.Type == tok {
			return p.currToken, nil
		}
	}
	err := p.parseError("expected token to be one of (%s), got '%s' instead",
		GetTokensString(tokenTypes), p.currToken.Type)
	return Token{}, err
}

func (p *Parser) nextTokenIn(tokenTypes []TokenID) bool {
	for _, TokenID := range tokenTypes {
		if p.nextToken.Type == TokenID {
			return true
		}
	}
	return false
}

func (p *Parser) currTokenIn(tokenTypes []TokenID) bool {
	for _, TokenID := range tokenTypes {
		if p.currToken.Type == TokenID {
			return true
		}
	}
	return false
}

func (p *Parser) requireToken(tok TokenID) error {
	if err := p.read(); err != nil {
		return err
	}
	if _, err := p.getExpectedToken(tok); err != nil {
		return err
	}
	return nil
}

func (p *Parser) requireTokenSequence(tokens []TokenID) error {
	for _, tok := range tokens {
		if err := p.requireToken(tok); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) registerUnaryExprFunction(TokenID TokenID, fn unaryExprFunction) {
	p.unaryExprFunctions[TokenID] = fn
}

func (p *Parser) registerBinExprFunction(TokenID TokenID, fn binExprFunctions) {
	p.binExprFunctions[TokenID] = fn
}

func (p *Parser) parseError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return errors.New(fmt.Sprintf("%s\nline:%d, pos %d", msg, p.currToken.Line, p.currToken.Col))
}
