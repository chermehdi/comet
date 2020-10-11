package eval

import (
	"github.com/chermehdi/comet/lexer"
	"github.com/chermehdi/comet/parser"
	"github.com/chermehdi/comet/std"
	"strings"
)

type Evaluator struct {
	Scope    *Scope
	Builtins map[string]*std.Builtin
}

// Constructs a new evaluator
// Each constructed evaluator has it's own Scope, i.e variables accessible from one Evaluator
// Are not accessible from another one.
func NewEvaluator() *Evaluator {
	ev := &Evaluator{
		Builtins: make(map[string]*std.Builtin),
		Scope:    NewScope(nil),
	}
	for _, builtin := range std.Builtins {
		ev.registerBuiltin(builtin)
	}
	return ev
}

type Scope struct {
	// The variables bound to this Scope instance
	Variables map[string]std.CometObject

	// The parent Scope if we are inside a function
	// if this is nil, this is the global Scope instance.
	Parent *Scope
}

// Creates a new Scope with the given parent.
func NewScope(parent *Scope) *Scope {
	store := make(map[string]std.CometObject)
	return &Scope{
		Variables: store,
		Parent:    parent,
	}
}

// Looks up the object bound to the varName
// The lookup should explore the parent(s) Scope as well ans should return a tuple (obj, true)
// if an object is bound to the given varName, and false otherwise.
func (sc *Scope) Lookup(varName string) (std.CometObject, bool) {
	obj, ok := sc.Variables[varName]
	if ok {
		return obj, ok
	}
	if sc.Parent != nil {
		return sc.Parent.Lookup(varName)
	}
	return obj, ok
}

// Stores the object and binds it to the given varName.
// The function will return true if this is overriding another variable with the same name in the current Scope.
func (sc *Scope) Store(varName string, obj std.CometObject) bool {
	_, has := sc.Variables[varName]
	sc.Variables[varName] = obj
	return has
}

func (sc *Scope) Clear(name string) {
	delete(sc.Variables, name)
}

// Evaluates the given node into a CometObject
// If the node is a statement a CometNop object is returned
// Errors are CometObject instances as well, and they are designed to block
// the evaluation process.
func (ev *Evaluator) Eval(node parser.Node) std.CometObject {
	switch n := node.(type) {
	case *parser.RootNode:
		return ev.evalRootNode(n.Statements)
	case *parser.PrefixExpression:
		return ev.evalPrefixExpression(n)
	case *parser.NumberLiteralExpression:
		return &std.CometInt{n.ActualValue}
	case *parser.BooleanLiteral:
		if n.ActualValue {
			return std.TrueObject
		} else {
			return std.FalseObject
		}
	case *parser.StringLiteral:
		return &std.CometStr{Value: n.Value, Size: len(n.Value)}
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
		return &std.CometReturnWrapper{result}
	case *parser.DeclarationStatement:
		return ev.evalDeclareStatement(n)
	case *parser.IdentifierExpression:
		return ev.evalIdentifier(n)
	case *parser.FunctionStatement:
		return ev.registerFunc(n)
	case *parser.CallExpression:
		result := ev.evalCallExpression(n)
		return unwrap(result)
	case *parser.AssignExpression:
		result := unwrap(ev.Eval(n.Value))
		ev.Scope.Store(n.VarName, result)
		return result
	case *parser.ForStatement:
		return unwrap(ev.evalForStatement(n))
	}
	return std.NopInstance
}

func unwrap(result std.CometObject) std.CometObject {
	if result.Type() == std.ReturnWrapper {
		unwrapped := result.(*std.CometReturnWrapper)
		return unwrapped.Value
	}
	return result
}

func (ev *Evaluator) evalRootNode(statements []parser.Statement) std.CometObject {
	var res std.CometObject = std.NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *std.CometReturnWrapper:
			return cur.Value
		case *std.CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalStatements(statements []parser.Statement) std.CometObject {
	var res std.CometObject = std.NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *std.CometReturnWrapper:
			return cur
		case *std.CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalPrefixExpression(n *parser.PrefixExpression) std.CometObject {
	res := ev.Eval(n.Right)
	if isError(res) {
		return res
	}
	switch n.Op.Type {
	case lexer.Minus:
		if res.Type() != std.IntType {
			return std.CreateError("Cannot apply operator (-) on none INTEGER type %s", res.Type())
		}
		result := res.(*std.CometInt)
		result.Value *= -1
		return result
	case lexer.Bang:
		if res.Type() != std.BoolType {
			return std.CreateError("Cannot apply operator (!) on none BOOLEAN type %s", res.Type())
		}
		result := res.(*std.CometBool)
		if result.Value {
			return std.FalseObject
		} else {
			return std.TrueObject
		}
	default:
		return std.CreateError("Unrecognized prefix operator %s", n.Op.Literal)
	}
}

func (ev *Evaluator) evalBinaryExpression(n *parser.BinaryExpression) std.CometObject {
	left := ev.Eval(n.Left)
	if isError(left) {
		return left
	}

	right := ev.Eval(n.Right)
	if isError(right) {
		return right
	}

	if left.Type() == std.IntType && right.Type() == std.IntType {
		return applyOp(n.Op.Type, left, right)
	}
	if left.Type() == std.BoolType && right.Type() == std.BoolType {
		return applyBoolOp(n.Op.Type, left, right)
	}
	if left.Type() == std.StrType && right.Type() == std.StrType {
		return applyStrOp(n.Op.Type, left, right)
	}
	if left.Type() == std.StrType || right.Type() == std.StrType {
		// one of the two is a string, the other one should be promoted to a string
		if n.Op.Type == lexer.Plus {
			return applyStrOp(n.Op.Type, std.ToString(left), std.ToString(right))
		} else if n.Op.Type == lexer.Mul && (left.Type() == std.IntType || right.Type() == std.IntType) {
			if left.Type() == std.IntType {
				leftValue := left.(*std.CometInt)
				rightValue := right.(*std.CometStr)
				return &std.CometStr{Value: strings.Repeat(rightValue.Value, int(leftValue.Value)), Size: int(leftValue.Value) * rightValue.Size}
			} else {
				leftValue := left.(*std.CometStr)
				rightValue := right.(*std.CometInt)
				return &std.CometStr{Value: strings.Repeat(leftValue.Value, int(rightValue.Value)), Size: int(rightValue.Value) * leftValue.Size}
			}
		} else {
			return std.CreateError("Cannot apply operation '%s' on operands of type '%s' and '%s'", n.Op.Literal, left.Type(), right.Type())
		}
	}
	if left.Type() != right.Type() {
		// operators == and != are applicable here, Objects with different types are always not equal in comet.
		switch n.Op.Type {
		case lexer.EQ:
			return std.FalseObject
		case lexer.NEQ:
			return std.TrueObject
		}
	}
	return std.CreateError("Cannot apply operator %s on given types %v and %v", n.Op.Literal, left.Type(), right.Type())
}

func (ev *Evaluator) evalConditional(n *parser.IfStatement) std.CometObject {
	predicateRes := ev.Eval(n.Test)
	if predicateRes.Type() != std.BoolType {
		return std.CreateError("Test part of the if statement should evaluate to CometBool, evaluated to %s instead", predicateRes.ToString())
	}
	result := predicateRes.(*std.CometBool)
	if result.Value {
		return ev.Eval(&n.Then)
	} else {
		return ev.Eval(&n.Else)
	}
}

func (ev *Evaluator) evalDeclareStatement(n *parser.DeclarationStatement) std.CometObject {
	value := ev.Eval(n.Expression)
	if isError(value) {
		return value
	}
	// TODO(chermehdi): add a shadowing diagnostic message if the store is overriding
	// an existing variable
	ev.Scope.Store(n.Identifier.Literal, value)
	return std.NopInstance
}

func (ev *Evaluator) evalIdentifier(n *parser.IdentifierExpression) std.CometObject {
	obj, found := ev.Scope.Lookup(n.Name)
	if !found {
		return std.CreateError("Identifier (%s) is not bounded to any value, have you tried declaring it?", n.Name)
	}
	return obj
}

func (ev *Evaluator) registerFunc(n *parser.FunctionStatement) std.CometObject {
	function := &std.CometFunc{
		Params: n.Parameters,
		Body:   n.Block,
	}
	ev.Scope.Store(n.Name, function)
	return function
}

func (ev *Evaluator) evalCallExpression(n *parser.CallExpression) std.CometObject {
	funcName := n.Name
	if ev.isBuiltinFunc(funcName) {
		args := make([]std.CometObject, 0)
		for i := range n.Arguments {
			args = append(args, ev.Eval(n.Arguments[i]))
		}
		return ev.invokeBuiltin(funcName, args...)
	}

	function, found := ev.Scope.Lookup(funcName)
	if !found {
		return std.CreateError("Cannot find callable symbol %s", funcName)
	}
	if function.Type() != std.FuncType {
		return std.CreateError("Cannot invoke none callable object of type %s", function.Type())
	}

	funObj, _ := function.(*std.CometFunc)
	callSiteScope := NewScope(ev.Scope)
	for i, param := range funObj.Params {
		callSiteScope.Variables[param.Name] = ev.Eval(n.Arguments[i])
	}
	oldScope := ev.Scope
	ev.Scope = callSiteScope
	result := ev.Eval(funObj.Body)
	ev.Scope = oldScope
	return result
}

func (ev *Evaluator) isBuiltinFunc(name string) bool {
	_, found := ev.Builtins[name]
	return found
}

func (ev *Evaluator) registerBuiltin(builtin *std.Builtin) {
	ev.Builtins[builtin.Name] = builtin
}

func (ev *Evaluator) invokeBuiltin(name string, args ...std.CometObject) std.CometObject {
	return ev.Builtins[name].Func(args...)
}

func (ev *Evaluator) evalForStatement(n *parser.ForStatement) std.CometObject {
	obj := ev.Eval(n.Range)
	switch obj.Type() {
	case std.RangeType:
		rangeObj := obj.(*std.CometRange)
		oldScope := ev.Scope
		curScope := NewScope(oldScope)
		ev.Scope = curScope
		for i := rangeObj.From.Value; i <= rangeObj.To.Value; i++ {
			ev.Scope.Store(n.Key.Name, &std.CometInt{Value: i})
			ev.Scope.Store(n.Value.Name, &std.CometInt{Value: i})
			ev.Eval(n.Body)
		}
		ev.Scope.Clear(n.Key.Name)
		ev.Scope.Clear(n.Value.Name)
		ev.Scope = oldScope
		return std.NopInstance
	default:
		panic("not implemented yet!!")
	}
}

func applyOp(op lexer.TokenType, left std.CometObject, right std.CometObject) std.CometObject {
	leftInt := left.(*std.CometInt)
	rightInt := right.(*std.CometInt)
	switch op {
	case lexer.Plus:
		return &std.CometInt{leftInt.Value + rightInt.Value}
	case lexer.Minus:
		return &std.CometInt{leftInt.Value - rightInt.Value}
	case lexer.Mul:
		return &std.CometInt{leftInt.Value * rightInt.Value}
	case lexer.Div:
		return &std.CometInt{leftInt.Value / rightInt.Value}
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
	case lexer.DotDot:
		return &std.CometRange{From: *leftInt, To: *rightInt}
	default:
		return std.CreateError("Cannot recognize binary operator %s", op)
	}
}

func applyStrOp(op lexer.TokenType, left std.CometObject, right std.CometObject) std.CometObject {
	leftStr := left.(*std.CometStr)
	rightStr := right.(*std.CometStr)
	switch op {
	case lexer.Plus:
		var sb strings.Builder
		sb.Grow(leftStr.Size + rightStr.Size)
		sb.WriteString(leftStr.Value)
		sb.WriteString(rightStr.Value)
		return &std.CometStr{Value: sb.String(), Size: leftStr.Size + rightStr.Size}
	default:
		return std.CreateError("Cannot execute binary operator '%s' on strings", op)
	}
}

func applyBoolOp(op lexer.TokenType, left std.CometObject, right std.CometObject) std.CometObject {
	leftInt := left.(*std.CometBool)
	rightInt := right.(*std.CometBool)
	switch op {
	case "==":
		return boolValue(leftInt.Value == rightInt.Value)
	case "!=":
		return boolValue(leftInt.Value != rightInt.Value)
	default:
		return std.CreateError("None-applicable operator %s for booleans", op)
	}
}

func boolValue(condition bool) *std.CometBool {
	if condition {
		return std.TrueObject
	}
	return std.FalseObject
}

func isError(obj std.CometObject) bool {
	return obj.Type() == std.ErrorType
}
