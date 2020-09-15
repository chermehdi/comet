package eval

import (
	"fmt"
	"github.com/chermehdi/comet/parser"
)

var (
	TrueObject  = &CometBool{true}
	FalseObject = &CometBool{false}
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
	}
	return nil
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
	switch n.Op.Literal {
	case "-":
		if res.Type() != IntType {
			panic(fmt.Sprintf("Cannot apply operator (-) on none integer type %s", res.ToString()))
		}
		result := res.(*CometInt)
		result.Value *= -1
		return result
	case "!":
		if res.Type() != BoolType {
			panic(fmt.Sprintf("Cannot apply operator (!) on none boolean type %s", res.ToString()))
		}
		result := res.(*CometBool)
		if result.Value {
			return FalseObject
		} else {
			return TrueObject
		}
	}
	return nil
}
