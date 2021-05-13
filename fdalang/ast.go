package fdalang

type AstNode interface {
	GetToken() Token
}

type AstExpression interface {
	AstNode
	Expression()
}

type AstStatement interface {
	AstNode
	Statement()
}

type AstStatementsBlock struct {
	Statements []AstStatement
}

type AstAssignment struct {
	Token Token
	Left  *AstIdentifier
	Value AstExpression
}

func (node *AstAssignment) Expression() {}

type AstStatementWithVoidedExpression struct {
	Token Token
	Expr  AstExpression
}

func (node *AstStatementWithVoidedExpression) Statement() {}

type AstStructFieldAssignment struct {
	Token Token
	Left  *AstStructFieldCall
	Value AstExpression
}

func (node *AstStructFieldAssignment) Expression() {}

type AstUnary struct {
	Token    Token
	Right    AstExpression
	Operator TokenID
}

func (node *AstUnary) Expression() {}

type AstEmptier struct {
	Token   Token
	Type    string
	IsArray bool
}

func (node *AstEmptier) Expression() {}

type AstBinOperation struct {
	Token    Token
	Left     AstExpression
	Right    AstExpression
	Operator TokenID
}

func (node *AstBinOperation) Expression() {}

type AstIdentifier struct {
	Token Token
	Value string
}

func (node *AstIdentifier) Expression() {}

type AstNumInt struct {
	Token Token
	Value int64
}

func (node *AstNumInt) Expression() {}

type AstNumFloat struct {
	Token Token
	Value float64
}

func (node *AstNumFloat) Expression() {}

type AstBoolean struct {
	Token Token
	Value bool
}

func (node *AstBoolean) Expression() {}

type AstArray struct {
	Token        Token
	ElementsType string
	Elements     []AstExpression
}

func (node *AstArray) Expression() {}

type AstArrayIndexCall struct {
	Token Token
	Left  AstExpression
	Index AstExpression
}

func (node *AstArrayIndexCall) Expression() {}

type AstReturn struct {
	Token       Token
	ReturnValue AstExpression
}

func (node *AstReturn) Statement() {}

type AstFunction struct {
	Token           Token
	Arguments       []*AstVarAndType
	ReturnType      string
	StatementsBlock *AstStatementsBlock
}

func (node *AstFunction) Expression() {}

type AstVarAndType struct {
	Token   Token
	VarType string
	Var     *AstIdentifier
}

type AstFunctionCall struct {
	Token     Token
	Function  AstExpression
	Arguments []AstExpression
}

func (node *AstFunctionCall) Expression() {}

type AstIfStatement struct {
	Token          Token
	Condition      AstExpression
	PositiveBranch *AstStatementsBlock
	ElseBranch     *AstStatementsBlock
}

func (node *AstIfStatement) Statement() {}

type AstEnumDefinition struct {
	Token    Token
	Name     string
	Elements []string
}

func (node *AstEnumDefinition) Statement() {}

type AstEnumElementCall struct {
	Token    Token
	EnumExpr AstExpression
	Element  *AstIdentifier
}

func (node *AstEnumElementCall) Expression() {}

type AstStructDefinition struct {
	Token  Token
	Name   string
	Fields map[string]*AstVarAndType
}

func (node *AstStructDefinition) Statement() {}

type AstStruct struct {
	Token  Token
	Ident  *AstIdentifier
	Fields []*AstAssignment
}

func (node *AstStruct) Expression() {}

type AstStructFieldCall struct {
	Token      Token
	StructExpr AstExpression
	Field      *AstIdentifier
}

func (node *AstStructFieldCall) Expression() {}

type AstCase struct {
	Token          Token
	Condition      AstExpression
	PositiveBranch *AstStatementsBlock
}

type AstSwitch struct {
	Token            Token
	Cases            []*AstCase
	SwitchExpression AstExpression
	DefaultBranch    *AstStatementsBlock
}

func (node *AstSwitch) Statement() {}

func (node *AstAssignment) GetToken() Token                    { return node.Token }
func (node *AstStructFieldAssignment) GetToken() Token         { return node.Token }
func (node *AstUnary) GetToken() Token                         { return node.Token }
func (node *AstBinOperation) GetToken() Token                  { return node.Token }
func (node *AstIdentifier) GetToken() Token                    { return node.Token }
func (node *AstNumInt) GetToken() Token                        { return node.Token }
func (node *AstNumFloat) GetToken() Token                      { return node.Token }
func (node *AstArray) GetToken() Token                         { return node.Token }
func (node *AstArrayIndexCall) GetToken() Token                { return node.Token }
func (node *AstBoolean) GetToken() Token                       { return node.Token }
func (node *AstReturn) GetToken() Token                        { return node.Token }
func (node *AstStatementWithVoidedExpression) GetToken() Token { return node.Token }
func (node *AstFunction) GetToken() Token                      { return node.Token }
func (node *AstVarAndType) GetToken() Token                    { return node.Token }
func (node *AstFunctionCall) GetToken() Token                  { return node.Token }
func (node *AstIfStatement) GetToken() Token                   { return node.Token }
func (node *AstStructDefinition) GetToken() Token              { return node.Token }
func (node *AstStruct) GetToken() Token                        { return node.Token }
func (node *AstStructFieldCall) GetToken() Token               { return node.Token }
func (node *AstEnumDefinition) GetToken() Token                { return node.Token }
func (node *AstEnumElementCall) GetToken() Token               { return node.Token }
func (node *AstSwitch) GetToken() Token                        { return node.Token }
func (node *AstEmptier) GetToken() Token                       { return node.Token }
func (node *AstStatementsBlock) GetToken() Token {
	if len(node.Statements) > 0 {
		return node.Statements[0].GetToken()
	}
	tok := Token{}
	return tok
}
