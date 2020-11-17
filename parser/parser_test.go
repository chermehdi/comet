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

func (t *TestingVisitor) VisitArrayAccess(access IndexAccess) {
	currentNode := t.expected[t.ptr]
	_, isIndexAccess := currentNode.(*IndexAccess)
	assert.True(t.t, isIndexAccess)
	t.ptr++
	access.Identifier.Accept(t)
	access.Index.Accept(t)
}

func (t *TestingVisitor) VisitExpression(Expression) {}

func (t *TestingVisitor) VisitStatement(Statement) {}

func (t *TestingVisitor) VisitStringLiteral(literal StringLiteral) {
	currentNode := t.expected[t.ptr]
	currentStringLiteral, ok := currentNode.(*StringLiteral)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentStringLiteral.Value, literal.Value)
	t.ptr++
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
	currentNode := t.expected[t.ptr]
	prefExpected, ok := currentNode.(*PrefixExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, prefExpected.Op.Literal, expression.Op.Literal)
	t.ptr++
	expression.Right.Accept(t)
}

func (t *TestingVisitor) VisitNumberLiteral(expression NumberLiteral) {
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

func (t *TestingVisitor) assertNumberLiteralNode(expression NumberLiteral) {
	currentNode := t.expected[t.ptr]
	currentNumberLiteral, ok := currentNode.(*NumberLiteral)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentNumberLiteral.ActualValue, expression.ActualValue)
	t.ptr++
}

func (t *TestingVisitor) VisitArrayLiteral(array ArrayLiteral) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*ArrayLiteral)
	assert.True(t.t, ok)
	t.ptr++
	for _, el := range array.Elements {
		el.Accept(t)
	}
}

func (t *TestingVisitor) VisitIdentifierExpression(expression IdentifierExpression) {
	currentNode := t.expected[t.ptr]
	currentIdentifierExpression, ok := currentNode.(*IdentifierExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentIdentifierExpression.Name, expression.Name)
	t.ptr++
}

func (t *TestingVisitor) VisitAssignExpression(assign AssignExpression) {
	currentNode := t.expected[t.ptr]
	currentAssignExpression, ok := currentNode.(*AssignExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, currentAssignExpression.VarName, assign.VarName)
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

func (t *TestingVisitor) VisitBlockStatement(statement BlockStatement) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*BlockStatement)
	assert.True(t.t, ok)
	t.ptr++
	for _, statement := range statement.Statements {
		statement.Accept(t)
	}
}

func (t *TestingVisitor) VisitIfStatement(statement IfStatement) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*IfStatement)
	assert.True(t.t, ok)
	t.ptr++
	statement.Test.Accept(t)
	statement.Then.Accept(t)
	statement.Else.Accept(t)
}

func (t *TestingVisitor) VisitForStatement(statement ForStatement) {
	currentNode := t.expected[t.ptr]
	_, ok := currentNode.(*ForStatement)
	assert.True(t.t, ok)
	t.ptr++
	statement.Key.Accept(t)
	statement.Value.Accept(t)
	statement.Range.Accept(t)
	statement.Body.Accept(t)
}

func (t *TestingVisitor) VisitFunctionStatement(statement FunctionStatement) {
	currentNode := t.expected[t.ptr]
	expectedFuncStatement, ok := currentNode.(*FunctionStatement)
	assert.True(t.t, ok)
	assert.Equal(t.t, expectedFuncStatement.Name, statement.Name)
	t.ptr++
	for _, parameter := range statement.Parameters {
		parameter.Accept(t)
	}
	statement.Block.Accept(t)
}

func (t *TestingVisitor) VisitCallExpression(expression CallExpression) {
	currentNode := t.expected[t.ptr]
	expectedCallExpression, ok := currentNode.(*CallExpression)
	assert.True(t.t, ok)
	assert.Equal(t.t, expectedCallExpression.Name, expression.Name)
	t.ptr++
	for _, arg := range expression.Arguments {
		arg.Accept(t)
	}
}

func (t *TestingVisitor) VisitBooleanLiteral(literal BooleanLiteral) {
	currentNode := t.expected[t.ptr]
	expectedBooleanLiteral, ok := currentNode.(*BooleanLiteral)
	assert.True(t.t, ok)
	assert.Equal(t.t, expectedBooleanLiteral.ActualValue, literal.ActualValue)
	t.ptr++
}

func TestParser_Parse_SimpleMathExpressions(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr:     "1",
			Expected: []Node{&NumberLiteral{ActualValue: int64(1)}},
		},
		{
			Expr: "1 + 21",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 - 21",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteral{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 * 21",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteral{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 / 21",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "/"}},
				&NumberLiteral{ActualValue: int64(21)},
			},
		},
		{
			Expr: "1 + 2 * 3 - 4",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteral{ActualValue: int64(3)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteral{ActualValue: int64(4)},
			},
		},
		{
			Expr: "(1)",
			Expected: []Node{
				&ParenthesisedExpression{},
				&NumberLiteral{ActualValue: int64(1)},
			},
		},
		{
			Expr: "1 * (1 + 2)",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&ParenthesisedExpression{},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
			},
		},
		{
			Expr: "1..2",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: ".."}},
				&NumberLiteral{ActualValue: int64(2)},
			},
		},
		{
			Expr: "1..a",
			Expected: []Node{
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: ".."}},
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
				&NumberLiteral{ActualValue: int64(1)},
			},
		},
		{
			Expr: "var a = 1 + 12",
			Expected: []Node{
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(12)},
			},
		},
		{
			Expr: `var a = 1 + 12
                   var b = (a + 12)`,
			Expected: []Node{
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(12)},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "b"},
				},
				&ParenthesisedExpression{},
				&IdentifierExpression{Name: "a"},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(12)},
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
				&NumberLiteral{ActualValue: int64(1)},
			},
		},
		{
			Expr: "return 1 * 2",
			Expected: []Node{
				&ReturnStatement{},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "*"}},
				&NumberLiteral{ActualValue: int64(2)},
			},
		},
		{
			Expr: "return 1 + (2 - 1)",
			Expected: []Node{
				&ReturnStatement{},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&ParenthesisedExpression{},
				&NumberLiteral{ActualValue: int64(2)},
				&BinaryExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteral{ActualValue: int64(1)},
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
				&NumberLiteral{ActualValue: int64(1)},
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

func TestParser_ParseBooleans(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: "true",
			Expected: []Node{
				&BooleanLiteral{ActualValue: true},
			},
		},
		{
			Expr: "false",
			Expected: []Node{
				&BooleanLiteral{ActualValue: false},
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

func TestParser_ParsePrefixExpression(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: "-1",
			Expected: []Node{
				&PrefixExpression{Op: lexer.Token{Literal: "-"}},
				&NumberLiteral{ActualValue: 1},
			},
		},
		{
			Expr: "!true",
			Expected: []Node{
				&PrefixExpression{Op: lexer.Token{Literal: "!"}},
				&BooleanLiteral{ActualValue: true},
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

func TestParser_ParsePrefixOperators(t *testing.T) {

	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `{
	var a = 1 + 2
	return a
}`,
			Expected: []Node{
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&ReturnStatement{},
				&IdentifierExpression{Name: "a"},
			},
		},
		{
			Expr: `{}`,
			Expected: []Node{
				&BlockStatement{},
			},
		},
		{
			Expr: `{}
			var a = 1 + 2`,
			Expected: []Node{
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
			},
		},
		{
			Expr: `{
				{
					var a = 1 + 2
				}
				{
					var b = 1 + 2
				}
				return a
			}`,
			Expected: []Node{
				&BlockStatement{},
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "b"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&ReturnStatement{},
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

func TestParser_ParseAssignExpression(t *testing.T) {

	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
	a = 10 + 1
`,
			Expected: []Node{
				&AssignExpression{VarName: "a"},
				&NumberLiteral{ActualValue: int64(10)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(1)},
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

func TestParser_ParseBlockStatement(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `{
	var a = 1 + 2
	return a
}`,
			Expected: []Node{
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&ReturnStatement{},
				&IdentifierExpression{Name: "a"},
			},
		},
		{
			Expr: `{}`,
			Expected: []Node{
				&BlockStatement{},
			},
		},
		{
			Expr: `{}
			var a = 1 + 2`,
			Expected: []Node{
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
			},
		},
		{
			Expr: `{
				{
					var a = 1 + 2
				}
				{
					var b = 1 + 2
				}
				return a
			}`,
			Expected: []Node{
				&BlockStatement{},
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{},
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "b"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&ReturnStatement{},
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

func TestParser_ParseIfStatement(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
				if a == 1 {
				}
`,
			Expected: []Node{
				&IfStatement{},
				&IdentifierExpression{Name: "a"},
				&BinaryExpression{Op: lexer.Token{Literal: "=="}},
				&NumberLiteral{ActualValue: int64(1)},
				&BlockStatement{},
				&BlockStatement{}, // accounting for the then empty block.
			},
		},
		{
			Expr: `
				if a == 1 {
					var a = 1 + 2	
				}else {}
`,
			Expected: []Node{
				&IfStatement{},
				&IdentifierExpression{Name: "a"},
				&BinaryExpression{Op: lexer.Token{Literal: "=="}},
				&NumberLiteral{ActualValue: int64(1)},
				&BlockStatement{}, // if close
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{}, // accounting for the then empty block.
			},
		},
		{
			Expr: `
				if (a == b) {
					var a = 1 + 2	
				}else {}
`,
			Expected: []Node{
				&IfStatement{},
				&ParenthesisedExpression{},
				&IdentifierExpression{Name: "a"},
				&BinaryExpression{Op: lexer.Token{Literal: "=="}},
				&IdentifierExpression{Name: "b"},
				&BlockStatement{}, // if close
				&DeclarationStatement{
					Identifier: lexer.Token{Literal: "a"},
				},
				&NumberLiteral{ActualValue: int64(1)},
				&BinaryExpression{Op: lexer.Token{Literal: "+"}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{}, // accounting for the then empty block.
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

func TestParser_ParseForStatement(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
						for k in 0..2 {
						}
		`,
			Expected: []Node{
				&ForStatement{},
				&IdentifierExpression{Name: "k"},
				// TODO: a placeholder identifier should be a singleton to avoid instance bloat.
				&IdentifierExpression{Name: "__empty__"},
				&NumberLiteral{ActualValue: int64(0)},
				&BinaryExpression{Op: lexer.Token{Literal: ".."}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{},
			},
		},
		{
			Expr: `
						for k, v in 0..2 {
						}
		`,
			Expected: []Node{
				&ForStatement{},
				&IdentifierExpression{Name: "k"},
				&IdentifierExpression{Name: "v"},
				&NumberLiteral{ActualValue: int64(0)},
				&BinaryExpression{Op: lexer.Token{Literal: ".."}},
				&NumberLiteral{ActualValue: int64(2)},
				&BlockStatement{},
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
func TestParser_Parse_ParseFunctionDeclaration(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
			func foo() {}
		`,
			Expected: []Node{
				&FunctionStatement{Name: "foo"},
				&BlockStatement{},
			},
		},
		{
			Expr: `
			func foo(a) {}
		`,
			Expected: []Node{
				&FunctionStatement{Name: "foo"},
				&IdentifierExpression{Name: "a"},
				&BlockStatement{},
			},
		},
		{
			Expr: `
			func foo(a) {}

			func bar() {}
		`,
			Expected: []Node{
				&FunctionStatement{Name: "foo"},
				&IdentifierExpression{Name: "a"},
				&BlockStatement{},
				&FunctionStatement{Name: "bar"},
				&BlockStatement{},
			},
		},
		{
			Expr: `
			func foo(a, b) {
				return 10
			}
		`,
			Expected: []Node{
				&FunctionStatement{Name: "foo"},
				&IdentifierExpression{Name: "a"},
				&IdentifierExpression{Name: "b"},
				&BlockStatement{},
				&ReturnStatement{},
				&NumberLiteral{10},
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

func TestParser_Parse_ShouldFailWrongFunctionCall(t *testing.T) {
	text := `func foo() {`
	parser := New(text)
	parser.Parse()
	assert.True(t, parser.Errors.HasAny())
}

func TestParser_Parse_ParseFunctionCall(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
			foo()
		`,
			Expected: []Node{
				&CallExpression{Name: "foo"},
				&BlockStatement{},
			},
		},
		{
			Expr: `
			foo(1 + 42, java, true)
		`,
			Expected: []Node{
				&CallExpression{Name: "foo"},
				&NumberLiteral{1},
				&BinaryExpression{Op: lexer.Token{Literal: lexer.Plus}},
				&NumberLiteral{42},
				&IdentifierExpression{"java"},
				&BooleanLiteral{
					ActualValue: true,
				},
			},
		},
		{
			Expr: `
			var result = foo(1 + 42, java, true)
		`,
			Expected: []Node{
				&DeclarationStatement{Identifier: lexer.Token{Literal: "result"}},
				&CallExpression{Name: "foo"},
				&NumberLiteral{1},
				&BinaryExpression{Op: lexer.Token{Literal: lexer.Plus}},
				&NumberLiteral{42},
				&IdentifierExpression{"java"},
				&BooleanLiteral{
					ActualValue: true,
				},
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

func TestParser_Parse_ParseArrayLiteral(t *testing.T) {
	tests := []struct {
		Expr     string
		Expected []Node
	}{
		{
			Expr: `
		 	[]
		 `,
			Expected: []Node{
				&ArrayLiteral{},
			},
		},
		{
			Expr: `
		 	[1, 2, 3]
		 `,
			Expected: []Node{
				&ArrayLiteral{},
				&NumberLiteral{ActualValue: 1},
				&NumberLiteral{ActualValue: 2},
				&NumberLiteral{ActualValue: 3},
			},
		},
		{
			Expr: `
			["comet"]
		`,
			Expected: []Node{
				&ArrayLiteral{},
				&StringLiteral{Value: "comet"},
			},
		},
		{
			Expr: `
			[[1, 2, 3], [42, 43, 44], [1]]
		`,
			Expected: []Node{
				&ArrayLiteral{},
				&ArrayLiteral{},
				&NumberLiteral{ActualValue: 1},
				&NumberLiteral{ActualValue: 2},
				&NumberLiteral{ActualValue: 3},

				&ArrayLiteral{},
				&NumberLiteral{ActualValue: 42},
				&NumberLiteral{ActualValue: 43},
				&NumberLiteral{ActualValue: 44},

				&ArrayLiteral{},
				&NumberLiteral{ActualValue: 1},
			},
		},
		{
			Expr: `
			a[0]	
		`,
			Expected: []Node{
				&IndexAccess{},
				&IdentifierExpression{Name: "a"},
				&NumberLiteral{ActualValue: 0},
			},
		},
		{
			Expr: `
			a[0] 
			1
		`,
			Expected: []Node{
				&IndexAccess{},
				&IdentifierExpression{Name: "a"},
				&NumberLiteral{ActualValue: 0},
				&NumberLiteral{ActualValue: 1},
			},
		},
		{
			Expr: `
			a["string"] 
			1
		`,
			Expected: []Node{
				&IndexAccess{},
				&IdentifierExpression{Name: "a"},
				&StringLiteral{Value: "string"},
				&NumberLiteral{ActualValue: 1},
			},
		},
		{
			Expr: `
			a[b()] 
			1
		`,
			Expected: []Node{
				&IndexAccess{},
				&IdentifierExpression{Name: "a"},
				&CallExpression{Name: "b"},
				&NumberLiteral{ActualValue: 1},
			},
		},
		{
			Expr: `
			a()[b()] 
			1
		`,
			Expected: []Node{
				&IndexAccess{},
				&CallExpression{Name: "a"},
				&CallExpression{Name: "b"},
				&NumberLiteral{ActualValue: 1},
			},
		},
	}

	for _, test := range tests {
		parser := New(test.Expr)
		rootNode := parser.Parse()
		assert.NotNil(t, rootNode)
		assert.False(t, parser.Errors.HasAny())

		testingVisitor := &TestingVisitor{
			expected: test.Expected,
			ptr:      0,
			t:        t,
		}
		rootNode.Accept(testingVisitor)
	}
}
