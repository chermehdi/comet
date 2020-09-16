package eval

import (
	"fmt"
	"github.com/chermehdi/comet/lexer"
	"github.com/chermehdi/comet/parser"
)

var (
	TrueObject  = &CometBool{true}
	FalseObject = &CometBool{false}
	NopInstance = &NopObject{}
)

type Evaluator struct {
	parser parser.Parser
}

func New() *Evaluator {
	return &Evaluator{}
}

func (ev *Evaluator) Eval(node parser.Node) CometObject {
	switch n := node.(type) {
	case *parser.RootNode:
		return ev.evalRootNode(n.Statements)
	case *parser.PrefixExpression:
		return ev.evalPrefixExpression(n)
	case *parser.NumberLiteralExpression:
		return &CometInt{n.ActualValue}
	case *parser.BooleanLiteral:
		if n.ActualValue {
			return TrueObject
		} else {
			return FalseObject
		}
	case *parser.BinaryExpression:
		return ev.evalBinaryExpression(n)
	case *parser.ParenthesisedExpression:
		return ev.Eval(n.Expression)
	case *parser.IfStatement:
		return ev.evalConditional(n)
	case *parser.BlockStatement:
		return ev.evalStatements(n.Statements)
	case *parser.ReturnStatement:
		result := ev.Eval(n.Expression)
		if isError(result) {
			return result
		}
		return &CometReturnWrapper{result}
	}
	return NopInstance
}

func (ev *Evaluator) evalRootNode(statements []parser.Statement) CometObject {
	var res CometObject = NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *CometReturnWrapper:
			return cur.Value
		case *CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalStatements(statements []parser.Statement) CometObject {
	var res CometObject = NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *CometReturnWrapper:
			return cur
		case *CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalPrefixExpression(n *parser.PrefixExpression) CometObject {
	res := ev.Eval(n.Right)
	if isError(res) {
		return res
	}
	switch n.Op.Type {
	case lexer.Minus:
		if res.Type() != IntType {
			return createError("Cannot apply operator (-) on none INTEGER type %s", res.Type())
		}
		result := res.(*CometInt)
		result.Value *= -1
		return result
	case lexer.Bang:
		if res.Type() != BoolType {
			return createError("Cannot apply operator (!) on none BOOLEAN type %s", res.Type())
		}
		result := res.(*CometBool)
		if result.Value {
			return FalseObject
		} else {
			return TrueObject
		}
	default:
		return createError("Unrecognized prefix operator %s", n.Op.Literal)
	}
}

func (ev *Evaluator) evalBinaryExpression(n *parser.BinaryExpression) CometObject {
	left := ev.Eval(n.Left)
	if isError(left) {
		return left
	}

	right := ev.Eval(n.Right)
	if isError(right) {
		return right
	}

	if left.Type() == IntType && right.Type() == IntType {
		return applyOp(n.Op.Type, left, right)
	}
	if left.Type() == BoolType && right.Type() == BoolType {
		return applyBoolOp(n.Op.Type, left, right)
	}
	if left.Type() != right.Type() {
		// operators == and != are applicable here, Objects with different types are always not equal in comet.
		switch n.Op.Type {
		case lexer.EQ:
			return FalseObject
		case lexer.NEQ:
			return TrueObject
		}
	}
	return createError("Cannot apply operator %s on given types %v and %v", n.Op.Literal, left.Type(), right.Type())
}

func (ev *Evaluator) evalConditional(n *parser.IfStatement) CometObject {
	predicateRes := ev.Eval(n.Test)
	if predicateRes.Type() != BoolType {
		return createError("Test part of the if statement should evaluate to CometBool, evaluated to %s instead", predicateRes.ToString())
	}
	result := predicateRes.(*CometBool)
	if result.Value {
		return ev.Eval(&n.Then)
	} else {
		return ev.Eval(&n.Else)
	}
}

func applyOp(op lexer.TokenType, left CometObject, right CometObject) CometObject {
	leftInt := left.(*CometInt)
	rightInt := right.(*CometInt)
	switch op {
	case lexer.Plus:
		return &CometInt{leftInt.Value + rightInt.Value}
	case lexer.Minus:
		return &CometInt{leftInt.Value - rightInt.Value}
	case lexer.Mul:
		return &CometInt{leftInt.Value * rightInt.Value}
	case lexer.Div:
		return &CometInt{leftInt.Value / rightInt.Value}
	case lexer.EQ:
		return boolValue(leftInt.Value == rightInt.Value)
	case lexer.NEQ:
		return boolValue(leftInt.Value != rightInt.Value)
	case lexer.LTE:
		return boolValue(leftInt.Value <= rightInt.Value)
	case lexer.LT:
		return boolValue(leftInt.Value < rightInt.Value)
	case lexer.GTE:
		return boolValue(leftInt.Value >= rightInt.Value)
	case lexer.GT:
		return boolValue(leftInt.Value > rightInt.Value)
	default:
		return createError("Cannot recognize binary operator %s", op)
	}
}

func createError(s string, params ...interface{}) CometObject {
	message := fmt.Sprintf(s, params...)
	return &CometError{
		message,
	}
}

func applyBoolOp(op lexer.TokenType, left CometObject, right CometObject) CometObject {
	leftInt := left.(*CometBool)
	rightInt := right.(*CometBool)
	switch op {
	case "==":
		return boolValue(leftInt.Value == rightInt.Value)
	case "!=":
		return boolValue(leftInt.Value != rightInt.Value)
	default:
		return createError("None-applicable operator %s for booleans", op)
	}
}

func boolValue(condition bool) *CometBool {
	if condition {
		return TrueObject
	}
	return FalseObject
}

func isError(obj CometObject) bool {
	return obj.Type() == ErrorType
}
