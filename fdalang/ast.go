package fdalang

type INode interface {
	GetToken() Token
}

type IExpression interface {
	INode
}

type IStatement interface {
	INode
}

type StatementsBlock struct {
	Statements []IStatement
}

type Assignment struct {
	Token Token
	Left  *Identifier
	Value IExpression
}

type StructFieldAssignment struct {
	Token Token
	Left  *StructFieldCall
	Value IExpression
}

type UnaryExpression struct {
	Token    Token
	Right    IExpression
	Operator string
}

type EmptierExpression struct {
	Token   Token
	Type    string
	IsArray bool
}

type BinExpression struct {
	Token    Token
	Left     IExpression
	Right    IExpression
	Operator string
}

type Identifier struct {
	Token Token
	Value string
}

type NumInt struct {
	Token Token
	Value int64
}

type NumFloat struct {
	Token Token
	Value float64
}

type Boolean struct {
	Token Token
	Value bool
}

type Array struct {
	Token        Token
	ElementsType string
	Elements     []IExpression
}

type ArrayIndexCall struct {
	Token Token
	Left  IExpression
	Index IExpression
}

type Return struct {
	Token       Token
	ReturnValue IExpression
}

type Function struct {
	Token           Token
	Arguments       []*VarAndType
	ReturnType      string
	StatementsBlock *StatementsBlock
}

type VarAndType struct {
	Token   Token
	VarType string
	Var     *Identifier
}

type FunctionCall struct {
	Token     Token
	Function  IExpression
	Arguments []IExpression
}

type IfStatement struct {
	Token          Token
	Condition      IExpression
	PositiveBranch *StatementsBlock
	ElseBranch     *StatementsBlock
}

type EnumDefinition struct {
	Token    Token
	Name     string
	Elements []string
}

type EnumElementCall struct {
	Token    Token
	EnumExpr IExpression
	Element  *Identifier
}

type StructDefinition struct {
	Token  Token
	Name   string
	Fields map[string]*VarAndType
}

type Struct struct {
	Token  Token
	Ident  *Identifier
	Fields []*Assignment
}

type StructFieldCall struct {
	Token      Token
	StructExpr IExpression
	Field      *Identifier
}

type Case struct {
	Token          Token
	Condition      IExpression
	PositiveBranch *StatementsBlock
}

type Switch struct {
	Token            Token
	Cases            []*Case
	SwitchExpression IExpression
	DefaultBranch    *StatementsBlock
}

func (node *Assignment) GetToken() Token            { return node.Token }
func (node *StructFieldAssignment) GetToken() Token { return node.Token }
func (node *UnaryExpression) GetToken() Token       { return node.Token }
func (node *BinExpression) GetToken() Token         { return node.Token }
func (node *Identifier) GetToken() Token            { return node.Token }
func (node *NumInt) GetToken() Token                { return node.Token }
func (node *NumFloat) GetToken() Token              { return node.Token }
func (node *Array) GetToken() Token                 { return node.Token }
func (node *ArrayIndexCall) GetToken() Token        { return node.Token }
func (node *Boolean) GetToken() Token               { return node.Token }
func (node *Return) GetToken() Token                { return node.Token }
func (node *Function) GetToken() Token              { return node.Token }
func (node *VarAndType) GetToken() Token            { return node.Token }
func (node *FunctionCall) GetToken() Token          { return node.Token }
func (node *IfStatement) GetToken() Token           { return node.Token }
func (node *StructDefinition) GetToken() Token      { return node.Token }
func (node *Struct) GetToken() Token                { return node.Token }
func (node *StructFieldCall) GetToken() Token       { return node.Token }
func (node *EnumDefinition) GetToken() Token        { return node.Token }
func (node *EnumElementCall) GetToken() Token       { return node.Token }
func (node *Case) GetToken() Token                  { return node.Token }
func (node *Switch) GetToken() Token                { return node.Token }
func (node *EmptierExpression) GetToken() Token     { return node.Token }
func (node *StatementsBlock) GetToken() Token {
	if len(node.Statements) > 0 {
		return node.Statements[0].GetToken()
	}
	tok := Token{}
	return tok
}
