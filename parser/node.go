package parser

import "github.com/chermehdi/comet/lexer"

type NodeVisitor interface {
	VisitExpression(Expression)
	VisitStatement(Statement)
	VisitRootNode(RootNode)
	VisitBinaryExpression(BinaryExpression)
	VisitPrefixExpression(PrefixExpression)
	VisitNumberLiteral(NumberLiteralExpression)
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
