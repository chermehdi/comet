package eval

import (
	lexer2 "github.com/chermehdi/comet/pkg/lexer"
	parser2 "github.com/chermehdi/comet/pkg/parser"
	std2 "github.com/chermehdi/comet/pkg/std"
	"strings"
)

type Evaluator struct {
	Scope    *Scope
	Builtins map[string]*std2.Builtin
	Types    map[string]*std2.CometStruct
}

// Param is a named parameter within the interpreter
type Param struct {
	Name string
	Val  std2.CometObject
}

// NewEvaluator Constructs a new evaluator
// Each constructed evaluator has it's own Scope, i.e variables accessible from one Evaluator
// Are not accessible from another one.
func NewEvaluator() *Evaluator {
	ev := &Evaluator{
		Builtins: make(map[string]*std2.Builtin),
		Types:    make(map[string]*std2.CometStruct),
		Scope:    NewScope(nil),
	}
	for _, builtin := range std2.Builtins {
		ev.registerBuiltin(builtin)
	}
	return ev
}

type Scope struct {
	// The variables bound to this Scope instance
	Variables map[string]std2.CometObject

	// The parent Scope if we are inside a function
	// if this is nil, this is the global Scope instance.
	Parent *Scope
}

// Creates a new Scope with the given parent.
func NewScope(parent *Scope) *Scope {
	store := make(map[string]std2.CometObject)
	return &Scope{
		Variables: store,
		Parent:    parent,
	}
}

// Looks up the object bound to the varName
// The lookup should explore the parent(s) Scope as well ans should return a tuple (obj, true)
// if an object is bound to the given varName, and false otherwise.
func (sc *Scope) Lookup(varName string) (std2.CometObject, bool) {
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
// The function will return true if the assignment of the variable has been done successfully
// returning false from this function implies that the variable has not been declared and should
// be handled appropriately.
func (sc *Scope) Store(varName string, obj std2.CometObject) bool {
	_, ok := sc.Variables[varName]
	if ok {
		sc.Variables[varName] = obj
		return true
	}
	if sc.Parent != nil {
		return sc.Parent.Store(varName, obj)
	}
	return false
}

// This function will create the symbol reference in the local scope.
func (sc *Scope) Declare(varName string, obj std2.CometObject) {
	sc.Variables[varName] = obj
}

func (sc *Scope) Clear(name string) {
	delete(sc.Variables, name)
}

// Evaluates the given node into a CometObject
// If the node is a statement a CometNop object is returned
// Errors are CometObject instances as well, and they are designed to block
// the evaluation process.
func (ev *Evaluator) Eval(node parser2.Node) std2.CometObject {
	switch n := node.(type) {
	case *parser2.RootNode:
		return ev.evalRootNode(n.Statements)
	case *parser2.PrefixExpression:
		return ev.evalPrefixExpression(n)
	case *parser2.NumberLiteral:
		return &std2.CometInt{Value: n.ActualValue}
	case *parser2.BooleanLiteral:
		if n.ActualValue {
			return std2.TrueObject
		} else {
			return std2.FalseObject
		}
	case *parser2.StringLiteral:
		return &std2.CometStr{Value: n.Value, Size: len(n.Value)}
	case *parser2.ArrayLiteral:
		return ev.evalArrayElements(n)
	case *parser2.BinaryExpression:
		return ev.evalBinaryExpression(n)
	case *parser2.ParenthesisedExpression:
		return ev.Eval(n.Expression)
	case *parser2.IfStatement:
		return ev.evalConditional(n)
	case *parser2.BlockStatement:
		return ev.evalStatements(n.Statements)
	case *parser2.ReturnStatement:
		result := ev.Eval(n.Expression)
		if isError(result) {
			return result
		}
		return &std2.CometReturnWrapper{Value: result}
	case *parser2.DeclarationStatement:
		return ev.evalDeclareStatement(n)
	case *parser2.IdentifierExpression:
		return ev.evalIdentifier(n)
	case *parser2.FunctionStatement:
		return ev.registerFunc(n)
	case *parser2.CallExpression:
		result := ev.evalCallExpression(n)
		return unwrap(result)
	case *parser2.AssignExpression:
		return ev.EvalAssignExpression(n)
	case *parser2.IndexAccess:
		return ev.evalArrayAccess(n)
	case *parser2.ForStatement:
		return unwrap(ev.evalForStatement(n))
	case *parser2.StructDeclarationStatement:
		return ev.evalStructDecl(n)
	case *parser2.NewCallExpr:
		return ev.evalNewCall(n)
	}
	return std2.NopInstance
}

func unwrap(result std2.CometObject) std2.CometObject {
	if result.Type() == std2.ReturnWrapper {
		unwrapped := result.(*std2.CometReturnWrapper)
		return unwrapped.Value
	}
	return result
}

func (ev *Evaluator) evalNewCall(expr *parser2.NewCallExpr) std2.CometObject {
	t, found := ev.Types[expr.Type]
	if !found {
		return std2.CreateError("Type '%s' not found", expr.Type)
	}
	instance := std2.NewInstance(t)
	params := make([]Param, len(expr.Args))
	constructor, found := t.GetConstructor()

	if !found {
		if len(expr.Args) > 0 {
			return std2.CreateError("Cannot find a defined constructor on the '%s' type, make sure to define an 'init' method on the struct", t.Name)
		}
		// it's okay to create a struct instance with no explicit constructor, if the constructor call didn't specify a parameter
		return instance
	}

	for i, param := range expr.Args {
		v := ev.Eval(param)
		if v.Type() == std2.ErrorType {
			return v
		}
		params[i] = Param{
			Name: constructor.Params[i].Name,
			Val:  v,
		}
	}
	res := ev.callOnObject("init", instance, params...)
	if res.Type() == std2.ErrorType {
		return res
	}
	return instance
}

func (ev *Evaluator) callOnObject(name string, object *std2.CometInstance, params ...Param) std2.CometObject {
	constructor, found := object.Struct.Methods[name]
	if found {
		callSiteScope := NewScope(ev.Scope)
		callSiteScope.Variables["this"] = object
		for _, p := range params {
			callSiteScope.Variables[p.Name] = p.Val
		}
		oldScope := ev.Scope
		ev.Scope = callSiteScope
		res := ev.Eval(constructor.Body)
		ev.Scope = oldScope
		return unwrap(res)
	}
	return std2.CreateError("Method '%s' Not found on instance of type '%s'", name, object.Struct.Name)
}

func (ev *Evaluator) evalRootNode(statements []parser2.Statement) std2.CometObject {
	var res std2.CometObject = std2.NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *std2.CometReturnWrapper:
			return cur.Value
		case *std2.CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalStructDecl(decl *parser2.StructDeclarationStatement) std2.CometObject {
	// The struct name should be defined in the global scope of the current
	// compilation unit.
	// The struct methods should be registered:
	//   - Scope the methods definitions with the struct declaration.
	//   - Register in the global scope with the a "cheeky naming scheme" -->
	//   Looks hacky
	s := &std2.CometStruct{Name: decl.Name, Methods: make(map[string]*std2.CometFunc, 0)}

	for _, m := range decl.Methods {
		fn := &std2.CometFunc{
			Name:   m.Name,
			Params: m.Parameters,
			Body:   m.Block,
		}
		if err := s.Add(fn); err != nil {
			return std2.CreateError(err.Error())
		}
	}

	ev.Types[s.Name] = s
	return std2.NopInstance
}

func (ev *Evaluator) evalStatements(statements []parser2.Statement) std2.CometObject {
	var res std2.CometObject = std2.NopInstance
	for _, st := range statements {
		res = ev.Eval(st)
		switch cur := res.(type) {
		case *std2.CometReturnWrapper:
			return cur
		case *std2.CometError:
			return cur
		}
	}
	return res
}

func (ev *Evaluator) evalPrefixExpression(n *parser2.PrefixExpression) std2.CometObject {
	res := ev.Eval(n.Right)
	if isError(res) {
		return res
	}
	switch n.Op.Type {
	case lexer2.Minus:
		if res.Type() != std2.IntType {
			return std2.CreateError("Cannot apply operator (-) on none INTEGER type %s", res.Type())
		}
		result := res.(*std2.CometInt)
		result.Value *= -1
		return result
	case lexer2.Bang:
		if res.Type() != std2.BoolType {
			return std2.CreateError("Cannot apply operator (!) on none BOOLEAN type %s", res.Type())
		}
		result := res.(*std2.CometBool)
		if result.Value {
			return std2.FalseObject
		} else {
			return std2.TrueObject
		}
	default:
		return std2.CreateError("Unrecognized prefix operator %s", n.Op.Literal)
	}
}

func (ev *Evaluator) SetField(instance *std2.CometInstance, name string, value parser2.Expression) {
	val := ev.Eval(value)
	instance.Fields[name] = val
}

func (ev *Evaluator) evalBinaryExpression(n *parser2.BinaryExpression) std2.CometObject {
	left := ev.Eval(n.Left)
	if isError(left) {
		return left
	}

	// Prioritize the dot operation
	if n.Op.Type == lexer2.Dot {
		as, ok := n.Right.(*parser2.AssignExpression)
		if ok {
			instance := left.(*std2.CometInstance)
			ev.SetField(instance, as.VarName, as.Value)
			return std2.NopInstance
		}
		id, ok := n.Right.(*parser2.IdentifierExpression)
		if ok {
			instance := left.(*std2.CometInstance)
			return instance.Fields[id.Name]
		}
		fn, ok := n.Right.(*parser2.CallExpression)
		if !ok {
			return std2.CreateError("Used '.' operator with none function element")
		}
		if left.Type() != std2.ObjType {
			// You can't call methods on none object types
			return std2.CreateError("Cannot call method '%s' on none object type", fn.Name)
		}

		instance := left.(*std2.CometInstance)
		method, found := instance.Struct.GetMethod(fn.Name)

		if !found {
			return std2.CreateError("Could not find method '%s' on type '%s'", fn.Name, instance.Struct.Name)
		}

		params := make([]Param, len(fn.Arguments))
		if len(method.Params) > len(fn.Arguments) {
			return std2.CreateError("Method '%s' on type '%s' expects at least %d parameters, %d were given",
				method.Name,
				instance.Struct.Name,
				len(method.Params),
				len(fn.Arguments))
		}

		for i, p := range method.Params {
			v := ev.Eval(fn.Arguments[i])
			if v.Type() == std2.ErrorType {
				return v
			}
			params[i] = Param{
				Name: p.Name,
				Val:  v,
			}
		}
		return ev.callOnObject(fn.Name, instance, params...)
	}

	right := ev.Eval(n.Right)
	if isError(right) {
		return right
	}

	if left.Type() == std2.IntType && right.Type() == std2.IntType {
		return applyOp(n.Op.Type, left, right)
	}
	if left.Type() == std2.BoolType && right.Type() == std2.BoolType {
		return applyBoolOp(n.Op.Type, left, right)
	}
	if left.Type() == std2.StrType && right.Type() == std2.StrType {
		return applyStrOp(n.Op.Type, left, right)
	}
	if left.Type() == std2.StrType || right.Type() == std2.StrType {
		// one of the two is a string, the other one should be promoted to a string
		if n.Op.Type == lexer2.Plus {
			return applyStrOp(n.Op.Type, std2.ToString(left), std2.ToString(right))
		} else if n.Op.Type == lexer2.Mul && (left.Type() == std2.IntType || right.Type() == std2.IntType) {
			if left.Type() == std2.IntType {
				leftValue := left.(*std2.CometInt)
				rightValue := right.(*std2.CometStr)
				return &std2.CometStr{Value: strings.Repeat(rightValue.Value, int(leftValue.Value)), Size: int(leftValue.Value) * rightValue.Size}
			} else {
				leftValue := left.(*std2.CometStr)
				rightValue := right.(*std2.CometInt)
				return &std2.CometStr{Value: strings.Repeat(leftValue.Value, int(rightValue.Value)), Size: int(rightValue.Value) * leftValue.Size}
			}
		} else {
			return std2.CreateError("Cannot apply operation '%s' on operands of type '%s' and '%s'", n.Op.Literal, left.Type(), right.Type())
		}
	}
	if left.Type() != right.Type() {
		// operators == and != are applicable here, Objects with different types are always not equal in comet.
		switch n.Op.Type {
		case lexer2.EQ:
			return std2.FalseObject
		case lexer2.NEQ:
			return std2.TrueObject
		}
	}
	return std2.CreateError("Cannot apply operator %s on given types %v and %v", n.Op.Literal, left.Type(), right.Type())
}

func (ev *Evaluator) evalConditional(n *parser2.IfStatement) std2.CometObject {
	predicateRes := ev.Eval(n.Test)
	if predicateRes.Type() != std2.BoolType {
		return std2.CreateError("Test part of the if statement should evaluate to CometBool, evaluated to %s instead", predicateRes.ToString())
	}
	result := predicateRes.(*std2.CometBool)
	if result.Value {
		return ev.Eval(&n.Then)
	} else {
		return ev.Eval(&n.Else)
	}
}

func (ev *Evaluator) evalDeclareStatement(n *parser2.DeclarationStatement) std2.CometObject {
	value := ev.Eval(n.Expression)
	if isError(value) {
		return value
	}
	// TODO(chermehdi): add a shadowing diagnostic message if the store is overriding
	// an existing variable
	ev.Scope.Declare(n.Identifier.Literal, value)
	return value
}

func (ev *Evaluator) evalIdentifier(n *parser2.IdentifierExpression) std2.CometObject {
	obj, found := ev.Scope.Lookup(n.Name)
	if !found {
		return std2.CreateError("Identifier (%s) is not bounded to any value, have you tried declaring it?", n.Name)
	}
	return obj
}

func (ev *Evaluator) registerFunc(n *parser2.FunctionStatement) std2.CometObject {
	function := &std2.CometFunc{
		Name:   n.Name,
		Params: n.Parameters,
		Body:   n.Block,
	}
	ev.Scope.Declare(n.Name, function)
	return function
}

func (ev *Evaluator) evalCallExpression(n *parser2.CallExpression) std2.CometObject {
	funcName := n.Name
	if ev.isBuiltinFunc(funcName) {
		args := make([]std2.CometObject, 0)
		for i := range n.Arguments {
			args = append(args, ev.Eval(n.Arguments[i]))
		}
		return ev.invokeBuiltin(funcName, args...)
	}

	function, found := ev.Scope.Lookup(funcName)
	if !found {
		return std2.CreateError("Cannot find callable symbol %s", funcName)
	}
	if function.Type() != std2.FuncType {
		return std2.CreateError("Cannot invoke none callable object of type %s", function.Type())
	}

	funObj, _ := function.(*std2.CometFunc)
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

func (ev *Evaluator) registerBuiltin(builtin *std2.Builtin) {
	ev.Builtins[builtin.Name] = builtin
}

func (ev *Evaluator) invokeBuiltin(name string, args ...std2.CometObject) std2.CometObject {
	return ev.Builtins[name].Func(args...)
}

func (ev *Evaluator) evalForStatement(n *parser2.ForStatement) std2.CometObject {
	obj := ev.Eval(n.Range)
	switch obj.Type() {
	case std2.RangeType:
		rangeObj := obj.(*std2.CometRange)
		oldScope := ev.Scope
		curScope := NewScope(oldScope)
		ev.Scope = curScope
		for i := rangeObj.From.Value; i <= rangeObj.To.Value; i++ {
			ev.Scope.Declare(n.Key.Name, &std2.CometInt{Value: i})
			ev.Scope.Declare(n.Value.Name, &std2.CometInt{Value: i})
			ev.Eval(n.Body)
		}
		ev.Scope.Clear(n.Key.Name)
		ev.Scope.Clear(n.Value.Name)
		ev.Scope = oldScope
		return std2.NopInstance
	default:
		panic("not implemented yet!!")
	}
}

func (ev *Evaluator) evalArrayElements(arr *parser2.ArrayLiteral) std2.CometObject {
	array := &std2.CometArray{
		Length: len(arr.Elements),
	}
	arrayContent := make([]std2.CometObject, array.Length)

	for i, expression := range arr.Elements {
		arrayContent[i] = ev.Eval(expression)
	}

	array.Values = arrayContent
	return array
}

func (ev *Evaluator) evalArrayAccess(arr *parser2.IndexAccess) std2.CometObject {
	array := ev.Eval(arr.Identifier)
	if array.Type() != std2.ArrayType {
		return std2.CreateError("Expected CometArray got %s", array.Type())
	}
	index := ev.Eval(arr.Index)
	if index.Type() != std2.IntType {
		return std2.CreateError("Expected CometInt got %s", index.Type())
	}
	indexVal := index.(*std2.CometInt)
	arrayVal := array.(*std2.CometArray)
	if indexVal.Value < 0 || indexVal.Value >= int64(arrayVal.Length) {
		return std2.CreateError("Array access out of bounds, array of length %d, index was: %d", arrayVal.Length, indexVal.Value)
	}
	return arrayVal.Values[int(indexVal.Value)]
}

func (ev *Evaluator) EvalAssignExpression(n *parser2.AssignExpression) std2.CometObject {
	_, found := ev.Scope.Lookup(n.VarName)
	if !found {
		return std2.CreateError("Identifier (%s) is not bounded to any value, have you tried declaring it?", n.VarName)
	}
	result := unwrap(ev.Eval(n.Value))
	ev.Scope.Store(n.VarName, result)
	return result
}

func applyOp(op lexer2.TokenType, left std2.CometObject, right std2.CometObject) std2.CometObject {
	leftInt := left.(*std2.CometInt)
	rightInt := right.(*std2.CometInt)
	switch op {
	case lexer2.Plus:
		return &std2.CometInt{Value: leftInt.Value + rightInt.Value}
	case lexer2.Minus:
		return &std2.CometInt{Value: leftInt.Value - rightInt.Value}
	case lexer2.Mul:
		return &std2.CometInt{Value: leftInt.Value * rightInt.Value}
	case lexer2.Div:
		return &std2.CometInt{Value: leftInt.Value / rightInt.Value}
	case lexer2.EQ:
		return boolValue(leftInt.Value == rightInt.Value)
	case lexer2.NEQ:
		return boolValue(leftInt.Value != rightInt.Value)
	case lexer2.LTE:
		return boolValue(leftInt.Value <= rightInt.Value)
	case lexer2.LT:
		return boolValue(leftInt.Value < rightInt.Value)
	case lexer2.GTE:
		return boolValue(leftInt.Value >= rightInt.Value)
	case lexer2.GT:
		return boolValue(leftInt.Value > rightInt.Value)
	case lexer2.DotDot:
		return &std2.CometRange{From: *leftInt, To: *rightInt}
	default:
		return std2.CreateError("Cannot recognize binary operator %s", op)
	}
}

func applyStrOp(op lexer2.TokenType, left std2.CometObject, right std2.CometObject) std2.CometObject {
	leftStr := left.(*std2.CometStr)
	rightStr := right.(*std2.CometStr)
	switch op {
	case lexer2.Plus:
		var sb strings.Builder
		sb.Grow(leftStr.Size + rightStr.Size)
		sb.WriteString(leftStr.Value)
		sb.WriteString(rightStr.Value)
		return &std2.CometStr{Value: sb.String(), Size: leftStr.Size + rightStr.Size}
	default:
		return std2.CreateError("Cannot execute binary operator '%s' on strings", op)
	}
}

func applyBoolOp(op lexer2.TokenType, left std2.CometObject, right std2.CometObject) std2.CometObject {
	leftInt := left.(*std2.CometBool)
	rightInt := right.(*std2.CometBool)
	switch op {
	case "==":
		return boolValue(leftInt.Value == rightInt.Value)
	case "!=":
		return boolValue(leftInt.Value != rightInt.Value)
	default:
		return std2.CreateError("None-applicable operator %s for booleans", op)
	}
}

func boolValue(condition bool) *std2.CometBool {
	if condition {
		return std2.TrueObject
	}
	return std2.FalseObject
}

func isError(obj std2.CometObject) bool {
	return obj.Type() == std2.ErrorType
}
