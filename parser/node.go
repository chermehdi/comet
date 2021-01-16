package parser

import (
	"fmt"
	"github.com/chermehdi/comet/lexer"
)

// API provided for all nodes types.
// Implementing a visitor will allow you to traverse the AST and perform some operation (printing, testing, code generation...)
// Without changing the Actual logic inside the AST.
//
// Example:
//    var visitor MyVisitor
//    rootNode.Accept(visitor)
//    visitor.getResult()
type NodeVisitor interface {
	VisitExpression(Expression)
	VisitStatement(Statement)

	VisitRootNode(RootNode)
	VisitBinaryExpression(BinaryExpression)
	VisitPrefixExpression(PrefixExpression)
	VisitNumberLiteral(NumberLiteral)
	VisitBooleanLiteral(BooleanLiteral)
	VisitStringLiteral(StringLiteral)
	VisitArrayLiteral(ArrayLiteral)
	VisitParenthesisedExpression(ParenthesisedExpression)
	VisitIdentifierExpression(IdentifierExpression)
	VisitCallExpression(CallExpression)
	VisitAssignExpression(AssignExpression)
	VisitArrayAccess(IndexAccess)
	VisitNewCall(NewCallExpr)

	VisitDeclarationStatement(DeclarationStatement)
	VisitReturnStatement(ReturnStatement)
	VisitBlockStatement(BlockStatement)
	VisitIfStatement(IfStatement)
	VisitFunctionStatement(FunctionStatement)
	VisitForStatement(ForStatement)
	VisitStructDeclaration(StructDeclarationStatement)
}

type Node interface {
	Literal() string
	Accept(NodeVisitor)
}

type Statement interface {
	Node
	Statement()
}

type Expression interface {
	Node
	Statement
	Expr()
}

type RootNode struct {
	Statements []Statement
}

func (r *RootNode) Accept(visitor NodeVisitor) {
	visitor.VisitRootNode(*r)
}

func (r *RootNode) Statement() {
	panic("implement me")
}

func (r *RootNode) Expr() {
	panic("implement me")
}

func (r *RootNode) Literal() string {
	return ""
}

type BinaryExpression struct {
	Op    lexer.Token
	Left  Expression
	Right Expression
}

func (e *BinaryExpression) Accept(visitor NodeVisitor) {
	visitor.VisitBinaryExpression(*e)
}

func (e *BinaryExpression) Literal() string {
	return e.Op.Literal
}

func (e *BinaryExpression) Statement() {
	panic("implement me")
}

func (e *BinaryExpression) Expr() {
	panic("implement me")
}

type PrefixExpression struct {
	Op    lexer.Token
	Right Expression
}

func (p *PrefixExpression) Accept(visitor NodeVisitor) {
	visitor.VisitPrefixExpression(*p)
}

func (p *PrefixExpression) Literal() string {
	return p.Op.Literal
}

func (p *PrefixExpression) Statement() {
	panic("implement me")
}

func (p *PrefixExpression) Expr() {
	panic("implement me")
}

type ParenthesisedExpression struct {
	Expression Expression
}

func (p *ParenthesisedExpression) Literal() string {
	panic("implement me")
}

func (p *ParenthesisedExpression) Accept(visitor NodeVisitor) {
	visitor.VisitParenthesisedExpression(*p)
}

func (p *ParenthesisedExpression) Statement() {
	panic("implement me")
}

func (p *ParenthesisedExpression) Expr() {
	panic("implement me")
}

type IdentifierExpression struct {
	Name string
}

func (i *IdentifierExpression) Literal() string {
	panic("implement me")
}

func (i *IdentifierExpression) Accept(visitor NodeVisitor) {
	visitor.VisitIdentifierExpression(*i)
}

func (i *IdentifierExpression) Statement() {
	panic("implement me")
}

func (i *IdentifierExpression) Expr() {
	panic("implement me")
}

type DeclarationStatement struct {
	varToken   lexer.Token
	Identifier lexer.Token
	Expression Expression
}

func (d *DeclarationStatement) Literal() string {
	panic("implement me")
}

func (d *DeclarationStatement) Accept(visitor NodeVisitor) {
	visitor.VisitDeclarationStatement(*d)
}

func (d *DeclarationStatement) Statement() {
	panic("implement me")
}

type ReturnStatement struct {
	returnToken lexer.Token
	Expression  Expression
}

func (d *ReturnStatement) Literal() string {
	panic("implement me")
}

func (d *ReturnStatement) Accept(visitor NodeVisitor) {
	visitor.VisitReturnStatement(*d)
}

func (d *ReturnStatement) Statement() {
	panic("implement me")
}

type BooleanLiteral struct {
	ActualValue bool
	Token       lexer.Token
}

func (b *BooleanLiteral) Literal() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) Accept(visitor NodeVisitor) {
	visitor.VisitBooleanLiteral(*b)
}

func (b *BooleanLiteral) Statement() {
	panic("implement me")
}

func (b *BooleanLiteral) Expr() {
	panic("implement me")
}

// Empty block for AST nodes when the block statement is optional
// This is instance is used to make the comparison easy.
var EmptyBlock = &BlockStatement{
	Statements: []Statement{},
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) Literal() string {
	return "BlockStatement"
}

func (b *BlockStatement) Accept(visitor NodeVisitor) {
	visitor.VisitBlockStatement(*b)
}

func (b *BlockStatement) Statement() {
	panic("implement me")
}

type IfStatement struct {
	Test Expression
	Then BlockStatement // this can be empty
	Else BlockStatement // this can be empty
}

func (i *IfStatement) Literal() string {
	return "IfStatement"
}

func (i *IfStatement) Accept(visitor NodeVisitor) {
	visitor.VisitIfStatement(*i)
}

func (i *IfStatement) Statement() {
	panic("implement me")
}

func newIfStatement() *IfStatement {
	return &IfStatement{
		Then: *EmptyBlock,
		Else: *EmptyBlock,
	}
}

type ForStatement struct {
	Key   *IdentifierExpression
	Value *IdentifierExpression
	Range Expression
	Body  *BlockStatement
}

func (f *ForStatement) Literal() string {
	panic("implement me")
}

func (f *ForStatement) Accept(visitor NodeVisitor) {
	visitor.VisitForStatement(*f)
}

func (f *ForStatement) Statement() {
	panic("implement me")
}

type FunctionStatement struct {
	Name       string
	Parameters []*IdentifierExpression
	Block      *BlockStatement
}

func (f *FunctionStatement) Literal() string {
	panic("Implement me!")
}

func (f *FunctionStatement) Accept(visitor NodeVisitor) {
	visitor.VisitFunctionStatement(*f)
}

func (f *FunctionStatement) Statement() {
	panic("implement me")
}

func newFunctionStatement() *FunctionStatement {
	return &FunctionStatement{
		Parameters: make([]*IdentifierExpression, 0),
		Block:      EmptyBlock,
	}
}

type CallExpression struct {
	Name      string
	Arguments []Expression
}

func (c *CallExpression) Literal() string {
	panic("implement me")
}

func (c *CallExpression) Accept(visitor NodeVisitor) {
	visitor.VisitCallExpression(*c)
}

func (c *CallExpression) Statement() {
	panic("implement me")
}

func (c *CallExpression) Expr() {
	panic("implement me")
}

type AssignExpression struct {
	VarName string
	Value   Expression
}

func (a *AssignExpression) Literal() string {
	panic("implement me")
}

func (a *AssignExpression) Accept(visitor NodeVisitor) {
	visitor.VisitAssignExpression(*a)
}

func (a *AssignExpression) Statement() {
	panic("implement me")
}

func (a *AssignExpression) Expr() {
	panic("implement me")
}

type NumberLiteral struct {
	ActualValue int64
}

func (n *NumberLiteral) Accept(visitor NodeVisitor) {
	visitor.VisitNumberLiteral(*n)
}

func (n *NumberLiteral) Literal() string {
	panic("implement me")
}

func (n *NumberLiteral) Statement() {
	panic("implement me")
}

func (n *NumberLiteral) Expr() {
	panic("implement me")
}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Literal() string {
	return s.Value
}

func (s *StringLiteral) Accept(visitor NodeVisitor) {
	visitor.VisitStringLiteral(*s)
}

func (s *StringLiteral) Statement() {
	panic("implement me")
}

func (s *StringLiteral) Expr() {
	panic("implement me")
}

type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) Literal() string {
	return fmt.Sprintf("ArrayLiteral(%d)", len(a.Elements))
}

func (a *ArrayLiteral) Accept(visitor NodeVisitor) {
	visitor.VisitArrayLiteral(*a)
}

func (a *ArrayLiteral) Statement() {
	panic("implement me")
}

func (a *ArrayLiteral) Expr() {
	panic("implement me")
}

type IndexAccess struct {
	Identifier Expression
	Index      Expression
}

func (i *IndexAccess) Literal() string {
	panic("implement me")
}

func (i *IndexAccess) Accept(visitor NodeVisitor) {
	visitor.VisitArrayAccess(*i)
}

func (i *IndexAccess) Statement() {
	panic("implement me")
}

func (i *IndexAccess) Expr() {
	panic("implement me")
}

type StructDeclarationStatement struct {
	Name    string
	Methods []*FunctionStatement
}

func (s *StructDeclarationStatement) Statement() {
	panic("implement me")
}

func (s *StructDeclarationStatement) Literal() string {
	panic("implement me")
}

func (s *StructDeclarationStatement) Accept(visitor NodeVisitor) {
	visitor.VisitStructDeclaration(*s)
}

type NewCallExpr struct {
	Type string
	Args []Expression
}

func (n *NewCallExpr) Expr() {
	panic("implement me!")
}

func (n *NewCallExpr) Statement() {
	panic("implement me!")
}

func (n *NewCallExpr) Literal() string {
	panic("implement me!")
}

func (n *NewCallExpr) Accept(visitor NodeVisitor) {
	visitor.VisitNewCall(*n)
}
