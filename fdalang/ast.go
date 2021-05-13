package fdalang

type Node interface {
	GetToken() Token
}

type Expression interface {
	Node
	Expression()
}

type Statement interface {
	Node
	Statement()
}

type StatementsBlock struct {
	Statements []Statement
}

type Assignment struct {
	Token Token
	Left  *Identifier
	Value Expression
}

func (node *Assignment) Expression() {}

type StatementWithVoidedExpression struct {
	Token Token
	Expr  Expression
}

func (node *StatementWithVoidedExpression) Statement() {}

type StructFieldAssignment struct {
	Token Token
	Left  *StructFieldCall
	Value Expression
}

func (node *StructFieldAssignment) Expression() {}

type UnaryExpression struct {
	Token    Token
	Right    Expression
	Operator string
}
func (node *UnaryExpression) Expression() {}

type EmptierExpression struct {
	Token   Token
	Type    string
	IsArray bool
}

func (node *EmptierExpression) Expression() {}

type BinExpression struct {
	Token    Token
	Left     Expression
	Right    Expression
	Operator string
}

func (node *BinExpression) Expression() {}

type Identifier struct {
	Token Token
	Value string
}

func (node *Identifier) Expression() {}

type NumInt struct {
	Token Token
	Value int64
}

func (node *NumInt) Expression() {}

type NumFloat struct {
	Token Token
	Value float64
}

func (node *NumFloat) Expression() {}

type Boolean struct {
	Token Token
	Value bool
}

func (node *Boolean) Expression() {}

type Array struct {
	Token        Token
	ElementsType string
	Elements     []Expression
}

func (node *Array) Expression() {}

type ArrayIndexCall struct {
	Token Token
	Left  Expression
	Index Expression
}

func (node *ArrayIndexCall) Expression() {}

type ReturnStatement struct {
	Token       Token
	ReturnValue Expression
}

func (node *ReturnStatement) Statement() {}

type Function struct {
	Token           Token
	Arguments       []*VarAndType
	ReturnType      string
	StatementsBlock *StatementsBlock
}

func (node *Function) Expression() {}

type VarAndType struct {
	Token   Token
	VarType string
	Var     *Identifier
}

type FunctionCall struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (node *FunctionCall) Expression() {}

type IfStatement struct {
	Token          Token
	Condition      Expression
	PositiveBranch *StatementsBlock
	ElseBranch     *StatementsBlock
}

func (node *IfStatement) Statement() {}

type EnumDefinition struct {
	Token    Token
	Name     string
	Elements []string
}

func (node *EnumDefinition) Statement() {}

type EnumElementCall struct {
	Token    Token
	EnumExpr Expression
	Element  *Identifier
}

func (node *EnumElementCall) Expression() {}

type StructDefinition struct {
	Token  Token
	Name   string
	Fields map[string]*VarAndType
}

func (node *StructDefinition) Statement() {}

type Struct struct {
	Token  Token
	Ident  *Identifier
	Fields []*Assignment
}

func (node *Struct) Expression() {}

type StructFieldCall struct {
	Token      Token
	StructExpr Expression
	Field      *Identifier
}

func (node *StructFieldCall) Expression() {}

type Case struct {
	Token          Token
	Condition      Expression
	PositiveBranch *StatementsBlock
}

type Switch struct {
	Token            Token
	Cases            []*Case
	SwitchExpression Expression
	DefaultBranch    *StatementsBlock
}

func (node *Switch) Statement() {}

func (node *Assignment) GetToken() Token                    { return node.Token }
func (node *StructFieldAssignment) GetToken() Token         { return node.Token }
func (node *UnaryExpression) GetToken() Token               { return node.Token }
func (node *BinExpression) GetToken() Token                 { return node.Token }
func (node *Identifier) GetToken() Token                    { return node.Token }
func (node *NumInt) GetToken() Token                        { return node.Token }
func (node *NumFloat) GetToken() Token                      { return node.Token }
func (node *Array) GetToken() Token                         { return node.Token }
func (node *ArrayIndexCall) GetToken() Token                { return node.Token }
func (node *Boolean) GetToken() Token                       { return node.Token }
func (node *ReturnStatement) GetToken() Token               { return node.Token }
func (node *StatementWithVoidedExpression) GetToken() Token { return node.Token }
func (node *Function) GetToken() Token                      { return node.Token }
func (node *VarAndType) GetToken() Token                    { return node.Token }
func (node *FunctionCall) GetToken() Token                  { return node.Token }
func (node *IfStatement) GetToken() Token                   { return node.Token }
func (node *StructDefinition) GetToken() Token              { return node.Token }
func (node *Struct) GetToken() Token                        { return node.Token }
func (node *StructFieldCall) GetToken() Token               { return node.Token }
func (node *EnumDefinition) GetToken() Token                { return node.Token }
func (node *EnumElementCall) GetToken() Token               { return node.Token }
func (node *Switch) GetToken() Token                        { return node.Token }
func (node *EmptierExpression) GetToken() Token             { return node.Token }
func (node *StatementsBlock) GetToken() Token {
	if len(node.Statements) > 0 {
		return node.Statements[0].GetToken()
	}
	tok := Token{}
	return tok
}