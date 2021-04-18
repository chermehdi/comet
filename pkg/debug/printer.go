package debug

import (
	"bytes"
	"fmt"
	parser2 "github.com/chermehdi/comet/pkg/parser"
)

const IndentWidth = 2

type PrintingVisitor struct {
	indent int
	buffer bytes.Buffer
}

func (p *PrintingVisitor) VisitNewCall(call parser2.NewCallExpr) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("NewCallExpression(%s)\n", call.Type))
	for _, ex := range call.Args {
		p.VisitExpression(ex)
	}
}

func (p *PrintingVisitor) VisitStructDeclaration(statement parser2.StructDeclarationStatement) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("StructDeclaration(Type=%s)\n", statement.Name))
	for _, method := range statement.Methods {
		p.VisitFunctionStatement(*method)
	}
}

func (p *PrintingVisitor) VisitArrayAccess(access parser2.IndexAccess) {
	p.printIndent()
	p.buffer.WriteString("IndexAccess")
	p.printIndent()
	p.VisitExpression(access.Identifier)
	p.VisitExpression(access.Index)
}

func (p *PrintingVisitor) VisitArrayLiteral(array parser2.ArrayLiteral) {
	p.printIndent()
	p.buffer.WriteString(array.Literal())
	for _, el := range array.Elements {
		p.VisitExpression(el)
	}
}

func (p *PrintingVisitor) VisitAssignExpression(expression parser2.AssignExpression) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("AssignmentExpression(%s)\n", expression.VarName))
}

func (p *PrintingVisitor) VisitForStatement(parser2.ForStatement) {
	panic("implement me")
}

func (p *PrintingVisitor) VisitStringLiteral(literal parser2.StringLiteral) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("StringLiteral(%s)\n", literal.Value))
}

func (p *PrintingVisitor) VisitIfStatement(statement parser2.IfStatement) {
	p.printIndent()
	p.buffer.WriteString("IfStatement\n")
	p.indent += IndentWidth
	statement.Test.Accept(p)
	p.buffer.WriteString("(Then)")
	statement.Then.Accept(p)
	p.buffer.WriteString("(Else)")
	statement.Else.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitBlockStatement(statement parser2.BlockStatement) {
	p.printIndent()
	p.buffer.WriteString("BlockStatement\n")
	p.indent += IndentWidth
	for _, statement := range statement.Statements {
		statement.Accept(p)
	}
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) printIndent() {
	for i := 0; i < p.indent; i++ {
		p.buffer.WriteRune(' ')
	}
}

func (p *PrintingVisitor) VisitExpression(parser2.Expression) {
	panic("implement me")
}

func (p *PrintingVisitor) VisitStatement(parser2.Statement) {
	panic("implement me")
}

func (p *PrintingVisitor) VisitRootNode(node parser2.RootNode) {
	p.printIndent()
	p.buffer.WriteString("RootNode\n")
	p.indent += IndentWidth
	for _, st := range node.Statements {
		st.Accept(p)
	}
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitBinaryExpression(expression parser2.BinaryExpression) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("Visiting a BinaryExpression (%s) \n", expression.Op.Literal))
	p.indent += IndentWidth
	expression.Left.Accept(p)
	expression.Right.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitPrefixExpression(expression parser2.PrefixExpression) {
	p.printIndent()
	p.buffer.WriteString("PrefixExpression\n")
	p.indent += IndentWidth
	expression.Right.Accept(p)
	expression.Right.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitNumberLiteral(expression parser2.NumberLiteral) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("Visiting a Number (%d)\n", expression.ActualValue))
}

func (p *PrintingVisitor) VisitParenthesisedExpression(expression parser2.ParenthesisedExpression) {
	p.printIndent()
	p.buffer.WriteString("ParenthesisedExpression\n")
	p.indent += IndentWidth
	expression.Expression.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) String() string {
	return p.buffer.String()
}

func (p *PrintingVisitor) VisitDeclarationStatement(statement parser2.DeclarationStatement) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("DeclarationStatement(%s)\n", statement.Identifier.Literal))
	p.indent += IndentWidth
	statement.Expression.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitIdentifierExpression(expression parser2.IdentifierExpression) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("IdentifierExpression(%s)\n", expression.Name))
}

func (p *PrintingVisitor) VisitReturnStatement(statement parser2.ReturnStatement) {
	p.printIndent()
	p.buffer.WriteString("ReturnStatement\n")
	p.indent += IndentWidth
	statement.Expression.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitBooleanLiteral(literal parser2.BooleanLiteral) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("BooleanLiteral (%v)\n", literal.ActualValue))
}

func (p *PrintingVisitor) VisitFunctionStatement(statement parser2.FunctionStatement) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("FuncStatement(Name='%s')\n", statement.Name))
	p.indent += IndentWidth
	p.printIndent()
	p.buffer.WriteString("Parameters: \n")
	for _, param := range statement.Parameters {
		param.Accept(p)
	}
	statement.Block.Accept(p)
	p.indent -= IndentWidth
}

func (p *PrintingVisitor) VisitCallExpression(expression parser2.CallExpression) {
	p.printIndent()
	p.buffer.WriteString(fmt.Sprintf("CallExpression(Name='%s')\n", expression.Name))
	p.indent += IndentWidth
	p.printIndent()
	p.buffer.WriteString("Parameters: \n")
	for _, arg := range expression.Arguments {
		arg.Accept(p)
	}
	p.indent -= IndentWidth
}
