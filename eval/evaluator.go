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
		return ev.evalStatements(n.Statements)
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
	}
	return NopInstance
}

func (ev *Evaluator) evalStatements(statements []parser.Statement) CometObject {
	var res CometObject
	for _, st := range statements {
		res = ev.Eval(st)
	}
	return res
}

func (ev *Evaluator) evalPrefixExpression(n *parser.PrefixExpression) CometObject {
	res := ev.Eval(n.Right)
	switch n.Op.Type {
	case lexer.Minus:
		if res.Type() != IntType {
			panic(fmt.Sprintf("Cannot apply operator (-) on none integer type %s", res.ToString()))
		}
		result := res.(*CometInt)
		result.Value *= -1
		return result
	case lexer.Bang:
		if res.Type() != BoolType {
			panic(fmt.Sprintf("Cannot apply operator (!) on none boolean type %s", res.ToString()))
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
	right := ev.Eval(n.Right)
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
