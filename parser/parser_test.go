package parser

import (
	"github.com/chermehdi/comet/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

// The testing visitor should assert on the structure of the tree while doing the traversal
// The expect nodes are given in in-order(ish) traversal.
// The change to the structure of the AST should make the TestingVisitor fail affected not updates tests.
type TestingVisitor struct {
	expected []Node
	ptr      int
	t        *testing.T
}

func (t *TestingVisitor) VisitExpression(expression Expression) {
}

func (t *TestingVisitor) VisitStatement(statement Statement) {
}

func (t *TestingVisitor) VisitRootNode(node RootNode) {
	for _, statement := range node.Statements {
		statement.Accept(t)
	}
}

func (t *TestingVisitor) VisitBinaryExpression(expression BinaryExpression) {
	expression.Left.Accept(t)
	t.assertBinaryExpression(expression)
	expression.Right.Accept(t)
}

func (t *TestingVisitor) VisitPrefixExpression(expression PrefixExpression) {
	panic("implement me")
}

func (t *TestingVisitor) VisitNumberLiteral(expression NumberLiteralExpression) {
	t.assertNumberLiteralNode(expression)
}

func (t *TestingVisitor) VisitParenthesisedExpression(expression ParenthesisedExpression) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*ParenthesisedExpression)
	assert.True(t.t, ok)
	t.ptr++
	expression.Expression.Accept(t)
}

func (t *TestingVisitor) assertBinaryExpression(expression BinaryExpression) {
	currentNode := t.expected[t.ptr]
	currentBinaryExpression, ok := currentNode.(*BinaryExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentBinaryExpression.Op.Literal, expression.Op.Literal)
	t.ptr++
}

func (t *TestingVisitor) assertNumberLiteralNode(expression NumberLiteralExpression) {
	currentNode := t.expected[t.ptr]
	currentNumberLiteral, ok := currentNode.(*NumberLiteralExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentNumberLiteral.ActualValue, expression.ActualValue)
	t.ptr++
}

func (t *TestingVisitor) VisitIdentifierExpression(expression IdentifierExpression) {
	currentNode := t.expected[t.ptr]
	currentIdentifierExpression, ok := currentNode.(*IdentifierExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentIdentifierExpression.Name, expression.Name)
	t.ptr++
}

func (t *TestingVisitor) VisitDeclarationStatement(statement DeclarationStatement) {
	currentNode := t.expected[t.ptr]
	currentDecStatement, ok := currentNode.(*DeclarationStatement)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentDecStatement.Identifier.Literal, statement.Identifier.Literal)
	t.ptr++
	statement.Expression.Accept(t)
}

func (t *TestingVisitor) VisitReturnStatement(statement ReturnStatement) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*ReturnStatement)
	assert.True(t.t, ok)
	t.ptr++
	statement.Expression.Accept(t)
}

func TestParser_Parse_SimpleMathExpressions(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr:     "1",
			Expected: []Node{&NumberLiteralExpression{ActualValue: int64(1)}},
		},
		{
			Expr: "1 + 21",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 - 21",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteralExpression{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 * 21",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteralExpression{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 / 21",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "/"}},
				&NumberLiteralExpression{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 + 2 * 3 - 4",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(2)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteralExpression{ActualValue: int64(3)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteralExpression{ActualValue: int64(4)},
			},
		},
		{
			Expr: "(1)",
			Expected: []Node{
				&ParenthesisedExpression{},
				&NumberLiteralExpression{ActualValue: int64(1)},
			},
		},
		{
			Expr: "1 * (1 + 2)",
			Expected: []Node{
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&ParenthesisedExpression{},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(2)},
			},
		},
	}
	for _, test := range tests {
		parser := New(test.Expr)
		rootNode := parser.Parse()
		assert.NotNil(t, rootNode)
		testingVisitor := &TestingVisitor{
			expected: test.Expected,
			ptr:      0,
			t:        t,
		}
		rootNode.Accept(testingVisitor)
	}
}

func TestParser_ParseDeclarationStatement(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: "var a = 1",
			Expected: []Node{
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteralExpression{ActualValue: int64(1)},
			},
		},
		{
			Expr: "var a = 1 + 12",
			Expected: []Node{
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(12)},
			},
		},
		{
			Expr: `var a = 1 + 12
                   var b = (a + 12)`,
			Expected: []Node{
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(12)},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "b"},
				},
				&ParenthesisedExpression{},
				&IdentifierExpression{Name: "a"},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteralExpression{ActualValue: int64(12)},
			},
		},
	}

	for _, test := range tests {
		parser := New(test.Expr)
		rootNode := parser.Parse()
		assert.NotNil(t, rootNode)
		testingVisitor := &TestingVisitor{
			expected: test.Expected,
			ptr:      0,
			t:        t,
		}
		rootNode.Accept(testingVisitor)
	}
}

func TestParser_ParseReturnStatement(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: "return 1",
			Expected: []Node{
				&ReturnStatement{},
				&NumberLiteralExpression{ActualValue: int64(1)},
			},
		},
		{
			Expr: "return 1 * 2",
			Expected: []Node{
				&ReturnStatement{},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteralExpression{ActualValue: int64(2)},
			},
		},
		{
			Expr: "return 1 + (2 - 1)",
			Expected: []Node{
				&ReturnStatement{},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&ParenthesisedExpression{},
				&NumberLiteralExpression{ActualValue: int64(2)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteralExpression{ActualValue: int64(1)},
			},
		},
		{
			Expr: "return value ",
			Expected: []Node{
				&ReturnStatement{},
				&IdentifierExpression{Name: "value"},
			},
		},
		{
			Expr: "return (1 + a)",
			Expected: []Node{
				&ReturnStatement{},
				&ParenthesisedExpression{},
				&NumberLiteralExpression{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&IdentifierExpression{Name: "a"},
			},
		},
	}

	for _, test := range tests {
		parser := New(test.Expr)
		rootNode := parser.Parse()
		assert.NotNil(t, rootNode)
		testingVisitor := &TestingVisitor{
			expected: test.Expected,
			ptr:      0,
			t:        t,
		}
		rootNode.Accept(testingVisitor)
	}
}
