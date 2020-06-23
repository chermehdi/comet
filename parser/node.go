package parser

import "github.com/chermehdi/comet/lexer"

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
	VisitNumberLiteral(NumberLiteralExpression)
	VisitBooleanLiteral(BooleanLiteral)
	VisitParenthesisedExpression(ParenthesisedExpression)
	VisitIdentifierExpression(IdentifierExpression)

	VisitDeclarationStatement(DeclarationStatement)
	VisitReturnStatement(statement ReturnStatement)
	VisitBlockStatement(statement BlockStatement)
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

type NumberLiteralExpression struct {
	ActualValue int64
}

func (n *NumberLiteralExpression) Accept(visitor NodeVisitor) {
	visitor.VisitNumberLiteral(*n)
}

func (n *NumberLiteralExpression) Literal() string {
	panic("implement me")
}

func (n *NumberLiteralExpression) Statement() {
	panic("implement me")
}

func (n *NumberLiteralExpression) Expr() {
	panic("implement me")
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
